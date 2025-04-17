import sys

ratios = {
    "composepost": 0.09979209979,
    "hometimeline": 0.7005428505,
    "poststorage": 0.9497747748,
    "socialgraph": 0.099792099790,
    "usertimeline": 0.4022002772
}

target = "hometimeline"
debug=True

file_name = sys.argv[1]
times = []
with open(f"results/{file_name}", 'r') as f:
    for line in f.readlines():
        line_output = line.strip()
        time_prefix = 'stop time: '
        if line_output.startswith(time_prefix):
            times.append(float(line_output[len(time_prefix):]))
baseline_throughput = 20000 * ratios[target] / (sum(times)/len(times)) 
# if debug: print(f"baseline_throughput: {baseline_throughput}")

base_processing_time = 2000000 / baseline_throughput
# if debug: print(f"base_processing_time: {base_processing_time}")
opt_t = opt_notcorrected = base_processing_time / 2
# if debug: print(f"opt_notcorrected: {opt_notcorrected}")
if len(sys.argv) > 2:
    premeasure_file = sys.argv[2]
    times = []
    with open(f"results/{premeasure_file}", 'r') as f:
        for line in f.readlines():
            line_output = line.strip()
            time_prefix = 'stop time: '
            if line_output.startswith(time_prefix):
                times.append(float(line_output[len(time_prefix):]))
    throughput = 20000 * ratios[target] / (sum(times)/len(times))
    # if debug: print(f"premesaured throughput: {throughput}")
    opt_processing_time =  2000000 / throughput
    # if debug: print(f"opt_processing_time: {opt_processing_time}")
    opt_t = opt_corrected = base_processing_time - opt_processing_time
    # if debug: print(f"opt_corrected: {opt_corrected}")

env = {}
env["SLOWPOKE_OPT_TIME"] = str(opt_t)

for service, ratio in ratios.items():
    if service == target:
        env[f"SLOWPOKE_PROCESSING_MICROS_{service.upper()}"] = 1000
        env[f"SLOWPOKE_IS_TARGET_SERVICE_{service.upper()}"] = "true"
        env[f"SLOWPOKE_POKER_BATCH_THRESHOLD_{service.upper()}"] = 100
        continue
    delay = opt_t * ratios[target] / (ratio * 2)
    env[f"SLOWPOKE_PROCESSING_MICROS_{service.upper()}"] = 0
    env[f"SLOWPOKE_DELAY_MICROS_{service.upper()}"] = str(delay)
    env[f"SLOWPOKE_IS_TARGET_SERVICE_{service.upper()}"] = "false"
    env[f"SLOWPOKE_POKER_BATCH_THRESHOLD_{service.upper()}"] = 100

if len(sys.argv) <= 3:
    for key, value in env.items():
        print(f"export {key}={value}")
    exit(0)

# if debug: print prediction
idx = int(sys.argv[3])
times = []
with open(f"results/groundtruth-{idx}.log", 'r') as f:
    for line in f.readlines():
        line_output = line.strip()
        time_prefix = 'stop time: '
        if line_output.startswith(time_prefix):
            times.append(float(line_output[len(time_prefix):]))
groundtruth = 20000 / (sum(times)/len(times))
if debug: print(f"groundtruth: {groundtruth}")

times = []
with open(f"results/slowdown-{idx}.log", 'r') as f:
    for line in f.readlines():
        line_output = line.strip()
        time_prefix = 'stop time: '
        if line_output.startswith(time_prefix):
            times.append(float(line_output[len(time_prefix):]))
throughput = 20000 / (sum(times)/len(times))
# if debug: print(f"slowdown: {throughput}")
predicted = 1000000 / (1000000 / throughput - (opt_notcorrected * ratios[target] / 2))
if debug: print(f"predicted: {predicted}")

times = []
with open(f"results/slowdown-corrected-{idx}.log", 'r') as f:
    for line in f.readlines():
        line_output = line.strip()
        time_prefix = 'stop time: '
        if line_output.startswith(time_prefix):
            times.append(float(line_output[len(time_prefix):]))
throughput = 20000 / (sum(times)/len(times))
# if debug: print(f"slowdown_corrected: {throughput}")
predicted_correction = 1000000 / (1000000 / throughput - (opt_corrected * ratios[target] / 2))
if debug: print(f"predicted_correction: {predicted_correction}")

# 4000000 / (3,158.2869254324 * 0.7005428505) = 1,807.8970816646
# baseline_throughput = 1,500.08378436712 * 0.7005428505
# baseline_processing_time = 2000000 / 1,050.8729702894 = 1,903.1796007173
# d = 1,903.1796007173 / 2 = 951.5898003586

# # delay = d * ratio[target] * cpu[service] / (ratio[service] * cpu[target])
# # actual_delay = delay / cpu[service] = d * ratio[target] / (ratio[service] * cpu[target])
# # compose = 951.5898003586 * 0.7005428505 / (0.09979209979 * 2 ) = 3,340.0912128955
# # poststorage = 951.5898003586 * 0.7005428505 / (0.9497747748 * 2 ) = 350.9407961431
# # socialgraph = 951.5898003586 * 0.7005428505 / (0.099792099790 * 2 ) = 3,340.0912128955
# # usertimeline = 951.5898003586 * 0.7005428505 / (0.4022002772 * 2 ) = 828.728209601
# # predicted = 1000000/ ( 1000000 / 1,506.6563418019 - (951.5898003586 * 0.7005428505 / 2) ) = 3,026.5734949887