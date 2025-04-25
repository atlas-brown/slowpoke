
file1="/home/ubuntu/mucache/slowpoke/synthetic/dag-cross-grpc-async/04-23-pokerpp-less-conn/dag-cross-grpc-async-service4-t8-c256-req20000-poker_batch_req100-n10-rep5.log"
file2="/home/ubuntu/mucache/slowpoke/synthetic/dag-cross-grpc-async/04-13-pokerpp-rm-deadlock-0maxconn/dag-cross-grpc-async-service4-t8-c512-req20000-poker_batch_req100-n10-rep5.log"



def parse_time(line):
    time = line.split(" ")[-1].strip()
    if "us" in time:
        return -1
    if "ms" in time:
        return float(time[:-2])
    if "s" in time:
        return float(time[:-1]) * 1000

def parse_log(file):
    all_data = []
    with open(file, 'r') as f:
        idx = 0
        lines = f.readlines()
        while idx < len(lines):
            line = lines[idx]
            if line.startswith("        50%  "):
                cur_data = []
                while idx < len(lines):
                    line = lines[idx]
                    for perc in [50, 75, 90, 99]:
                        if line.startswith(f"        {perc}%  "):
                            time = parse_time(line)
                            if time == -1:
                                break   
                            cur_data.append(time)
                    idx += 1
                    if line.startswith("    [exp] Throughput: "):
                        cur_data.append(float(line.split(" ")[-1].strip()))
                        break
                # print(cur_data)
                if len(cur_data) == 9:
                    all_data.append(cur_data.copy())
            idx += 1
    return all_data

data1 = parse_log(file1)
print(len(data1))
data2 = parse_log(file2)
print(len(data2))

for i in range(len(data1)):
    const1 = data1[i][7]/(3.8)
    const2 = data2[i][7]/(3.8)
    print(f"const1: {const1}, const2: {const2}, throughput1: {data1[i][8]}, throughput2: {data2[i][8]}")