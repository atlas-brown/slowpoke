import os

DIR = "/home/ubuntu/mucache/slowpoke/synthetic/two-services/interference-microbenchmark"
map_multiplier = {
    "combined": 1000
}
map_ = {}
for file in sorted(os.listdir(DIR)):
    service, delay, poker = file.split(".")[0].split("-")
    map_[service] = map_.get(service, {})
    map_[service][delay] = map_.get(delay, {})
    map_[service][delay][poker] = {}
    with open(os.path.join(DIR, file), "r") as f:
        # Process the file
        times = []
        for line in f.readlines():
            time_prefix = 'stop time: '
            if line.startswith(time_prefix):
                times.append(float(line[len(time_prefix):]))
        throughput = 40000 / (sum(times) / len(times))
        processing_time = 2000000 / throughput
        map_[service][delay][poker]["times"] = times
        map_[service][delay][poker]["throughput"] = throughput
        map_[service][delay][poker]["processing_time"] = processing_time
        # print(f"Service: {service}, Delay: {delay}, Poker: {poker}, Processing Time: {processing_time:.2f}, Throughput: {throughput:.2f}")
        poker_batch = int(poker.split("r")[1])
        if service == "combined":
            print(file)
            print(f"Service: {service}, Delay: {delay}, Poker: {poker}, Processing Time: {processing_time:.2f}, Throughput: {throughput:.2f}, cal: {poker_batch*throughput/10e6:.2f}")
            

# base = map_["base"]["base"]["processing_time"]
# print(f"Base Processing Time: {base:.2f}")
# # Print the results
# for service in map_:
#     base = map_[service]["base"]["processing_time"]
#     print(f"Service: {service}, base processing time {base:.2f}")
#     for delay in map_[service]:
#         if delay == "base":
#             continue
#         delay_int = int(delay.split("d")[1])
#         expected = base + delay_int*2
#         print(f"    Delay: {delay}, expected processing time {expected:.2f}")
#         for poker in map_[service][delay]:
#             print(f"      {poker} - Processing Time: {map_[service][delay][poker]['processing_time']:.2f} && Diff: {expected - map_[service][delay][poker]['processing_time']:.2f}")
