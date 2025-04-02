import os

map_multiplier = {
    "combined": 1000
}
map_ = {}

# DIR = "/home/ubuntu/mucache/slowpoke/synthetic/two-services/interference-microbenchmark"
# for file in sorted(os.listdir(DIR)):
#     service, delay, poker = file.split(".")[0].split("-")
#     map_[service] = map_.get(service, {})
#     map_[service][delay] = map_.get(delay, {})
#     map_[service][delay][poker] = {}
#     with open(os.path.join(DIR, file), "r") as f:
#         # Process the file
#         times = []
#         for line in f.readlines():
#             time_prefix = 'stop time: '
#             if line.startswith(time_prefix):
#                 times.append(float(line[len(time_prefix):]))
#         throughput = 40000 / (sum(times) / len(times))
#         processing_time = 2000000 / throughput
#         map_[service][delay][poker]["times"] = times
#         map_[service][delay][poker]["throughput"] = throughput
#         map_[service][delay][poker]["processing_time"] = processing_time
#         # print(f"Service: {service}, Delay: {delay}, Poker: {poker}, Processing Time: {processing_time:.2f}, Throughput: {throughput:.2f}")
#         poker_batch = int(poker.split("r")[1])
#         if service == "combined":
#             print(file)
#             print(f"Service: {service}, Delay: {delay}, Poker: {poker}, Processing Time: {processing_time:.2f}, Throughput: {throughput:.2f}, cal: {poker_batch*throughput/10e6:.2f}")
            

# DIR = "/home/ubuntu/mucache/slowpoke/synthetic/two-services/locker-microbenchmark"
DIR = "/home/ubuntu/mucache/slowpoke/synthetic/two-services/move-timer-read-microbenchmark"
for file in sorted(os.listdir(DIR)):
    delay, poker = file.split(".")[0].split("-")
    map_[delay] = map_.get(delay, {})
    map_[delay][poker] = map_[delay].get(poker, {})
    map_[delay][poker] = {}
    with open(os.path.join(DIR, file), "r") as f:
        # Process the file
        times = []
        for line in f.readlines():
            time_prefix = 'stop time: '
            if line.startswith(time_prefix):
                times.append(float(line[len(time_prefix):]))
        if len(times) == 0:
            del map_[delay][poker]
            continue
        throughput = 40000 / (sum(times) / len(times))
        processing_time = 2000000 / throughput
        map_[delay][poker]["times"] = times
        map_[delay][poker]["throughput"] = throughput
        map_[delay][poker]["processing_time"] = processing_time

# try:
#     base = map_["base"]["base"]["processing_time"]
#     print(f"Base Processing Time: {base:.2f}")
# except:
#     pass
# Print the results
for delay in map_:
    print(f"Exp: {delay}")
    for poker in map_[delay]:
        processing_time = map_[delay][poker]["processing_time"]
        throughput = map_[delay][poker]["throughput"]
        times = map_[delay][poker]["times"]
        # print(f"    {poker}: Time: {processing_time:.2f}, Tput: {throughput:.2f}, Diff: {(processing_time - expected):.2f}")
        print(f"    {poker}: Time: {processing_time:.2f}, Tput: {throughput:.2f}")
        # print(f"Service: {service}, Delay: {delay}, Poker: {poker}, Processing Time: {processing_time:.2f}, Throughput: {throughput:.2f}")
