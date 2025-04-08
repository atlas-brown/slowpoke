#!python3

import os
import subprocess
import copy
import argparse
import config

class Runner:
    def __init__(self, args):
        self.benchmark = args.benchmark
        self.num_threads = args.num_threads
        self.num_conns = args.num_conns
        self.target_service = args.target_service
        self.request_type = args.request_type
        self.repetitions = args.repetitions
        self.target_num_exp = args.num_exp
        # self.duration = args.duration
        self.pre_run = args.pre_run
        self.num_req = args.num_req
        self.groundtruths = []
        self.slowdowns = []
        self.predicteds = []
        self.errs = []
        self.poker_batch_req = args.poker_batch_req
        # self.poker_relative_batch = args.poker_relative_batch
        if args.clien_cpu_quota != -1:
            self.client_cpu_quota = args.clien_cpu_quota
        else:
            self.client_cpu_quota = 2 if self.benchmark == "synthetic" else 4
        self.random_seed = args.random_seed
        self.request_ratio = config.get_request_ratio(self.benchmark, self.request_type)
        self.baseline_service_processing_time = config.get_baseline_service_processing_time(self.benchmark, self.request_type, self.target_service, self.random_seed)
        self.cpu_quota = config.get_cpu_quota(self.benchmark, self.request_type)
        self.target_processing_time_range = [0, self.baseline_service_processing_time[self.target_service]]
        self.baseline_throughputs = []
        self.poker_batch = args.poker_batch
    
    def get_env_for_print(self, env):
        env_p = {}
        for key, value in env.items():
            if key.startswith("SLOWPOKE") or key.startswith("PROCESSING_TIME") or key.startswith("CLIENT"):
                env_p[key] = value
        return env_p
    
    def exp(self, service_delay, processing_time):
        env = os.environ.copy()
        env["CLIENT_CPU_QUOTA"] = str(self.client_cpu_quota)
        env["SLOWPOKE_POKER_BATCH_THRESHOLD"] = str(self.poker_batch)
        if self.benchmark == "synthetic":
            # This is used to set the processing time for synthetic benchmarks
            for service, processing_time in processing_time.items():
                env[f"PROCESSING_TIME_{service.upper()}"] = str(round(processing_time/1e6, 8))
        else:
            for service, p_ in processing_time.items():
                env[f"SLOWPOKE_PROCESSING_MICROS_{service.upper()}"] = str(p_)
        # expected_req_num = -1
        # for service, delay in service_delay.items():
        #     # we need to make sure the batch size for any service is larger than self.poker_relative_batch
        #     cur_service_expected_num = self.poker_relative_batch / delay # num of calls needed for this service
        #     expected_req_num = max(cur_service_expected_num/self.request_ratio[service], expected_req_num)
        for service, delay in service_delay.items():
            d = delay/self.cpu_quota[service]
            env[f"SLOWPOKE_DELAY_MICROS_{service.upper()}"] = str(d)
            # batch_size = max(10, int(self.poker_batch_req * 1000 * self.request_ratio[service] * d - 10))
            # env[f"SLOWPOKE_POKER_BATCH_THRESHOLD_{service.upper()}"] = str(batch_size)
            env[f"SLOWPOKE_POKER_BATCH_THRESHOLD_{service.upper()}"] = str(self.poker_batch_req)
            if service == self.target_service:
                env[f"SLOWPOKE_IS_TARGET_SERVICE_{service.upper()}"] = "true"
            else:
                env[f"SLOWPOKE_IS_TARGET_SERVICE_{service.upper()}"] = "false"
        if self.pre_run:
            env["SLOWPOKE_PRERUN"] = "true" # Disable request counting during normal execution
        cmd = f"bash run.sh {self.benchmark} {self.request_type} {self.num_threads} {self.num_conns} {self.num_req}"
        print(f"    [exp] Running (pre_run: {self.pre_run}) workload {self.benchmark}/{self.request_type} request with the following configuration: {self.get_env_for_print(env)}", flush=True)
        process = subprocess.Popen(cmd, shell=True, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        print(f"    [exp] Executing cmd `{cmd}`:")
        throughput = -1
        times = []
        for line in process.stdout:
            line_output = line.decode().strip()
            print(f"        {line_output}", flush=True)
            time_prefix = 'stop time: '
            if line_output.startswith(time_prefix):
                times.append(float(line_output[len(time_prefix):]))
        if process.stderr:
            print(f"        > STDERR", flush=True)
            for line in process.stderr:
                line = line.decode().strip()
                if "Dload" in line or "% Total" in line or "--:--:--" in line:
                    continue
                print(f"        > {line}", flush=True)
        if process.wait() != 0:
            print(f"    [exp] Error running {cmd}")
            return 0
        throughput = self.num_req / (sum(times)/len(times))
        print(f"    [exp] Times: {times}")
        print(f"    [exp] Throughput: {throughput}")
        return throughput

    def run(self):
        service_delay = copy.deepcopy(self.baseline_service_processing_time)
        for service in service_delay:
            service_delay[service] = 0
        processing_time = copy.deepcopy(self.baseline_service_processing_time)
        processing_time_diff = self.target_processing_time_range[1]-self.target_processing_time_range[0]
        processing_time_range = list(range(self.target_processing_time_range[0], int(self.target_processing_time_range[1]), int(processing_time_diff//self.target_num_exp)))[:self.target_num_exp]
        print(f"[test.py] Actual processing time range: {processing_time_range}")

        # baseline
        print(f"[test.py] Running baseline experiment", flush=True)
        baseline_throughput = self.exp(service_delay, processing_time)
        while baseline_throughput == 0:
            print("[test.py] Found 0 throughput, rerun experiment")
            baseline_throughput = self.exp(service_delay, processing_time)
        print(f"[test.py] Baseline throughput: {baseline_throughput}", flush=True)
        self.baseline_throughputs.append(baseline_throughput)

        slowdown = []
        groundtruth = []
        predicted = []
        err = []
        processing_time[self.target_service] = self.baseline_service_processing_time[self.target_service]
        for i, p_t in enumerate(processing_time_range):
            print(f"[test.py] Running {i}th optmization experiment", flush=True)

            print(f"[test.py] Running {i}th groundtruth exp", flush=True)
            processing_time[self.target_service] = p_t
            for service in service_delay:
                service_delay[service] = 0
            res = self.exp(service_delay, processing_time)
            while int(res) == 0:
                print("[test.py] Found 0 throughtput, rerun experiment")
                res = self.exp(service_delay, processing_time)
            groundtruth.append(res)

            print(f"[test.py] Running {i}th slowdown exp", flush=True)
            processing_time[self.target_service] = self.baseline_service_processing_time[self.target_service]
            for service in service_delay:
                if service != self.target_service:
                    if self.request_ratio[service] == 0:
                        delay = 0
                    else:
                        delay = int(
                            ((self.target_processing_time_range[1] - p_t)*self.request_ratio[self.target_service])*self.cpu_quota[service] / (self.request_ratio[service]*self.cpu_quota[self.target_service])
                        )
                    service_delay[service] = delay
            res = self.exp(service_delay, processing_time)
            while int(res) == 0:
                print("[test.py] Found 0 throughput, rerun experiment")
                res = self.exp(service_delay, processing_time)
            slowdown.append(res)

            try:
                delay = (self.target_processing_time_range[1] - p_t)*self.request_ratio[self.target_service]/self.cpu_quota[self.target_service]
                predicted_throughput = 1000000/(1000000/slowdown[-1]-delay)
            except:
                print("[test.py] Error: Division by zero")
                predicted_throughput = -1
            predicted.append(predicted_throughput)
            err.append((predicted[-1]-groundtruth[-1])*100/groundtruth[-1])
            print(f"[test.py] Finished running {i}th optmization experiment: groundtruth->{groundtruth[-1]}, slowdown->{slowdown[-1]}, predicted->{predicted[-1]}, err->{err[-1]}", flush=True)
        
        print("[test.py] Baseline throughput: ", baseline_throughput, flush=True)
        print("[test.py] Groundtruth: ", groundtruth, flush=True)
        print("[test.py] Slowdown: ", slowdown, flush=True)
        print("[test.py] Predicted: ", predicted, flush=True)
        print("[test.py] Error percentage: ", err, flush=True)

        self.groundtruths.append(groundtruth)
        self.slowdowns.append(slowdown)
        self.predicteds.append(predicted)
        self.errs.append(err)
    
    def start(self):
        print("[test.py] Starting experiment...", flush=True)
        for i in range(self.repetitions):
            print(f"[test.py] >>>>>>>>>>>>>>>>>>>>>>>>>>>>> Running experiment {i}...")
            self.run()
    
    def print_summary(self):
        print("[test.py] ++++++++++++++++++++++++++++++++ Summary: ")
        for i in range(len(self.groundtruths)):
            print(f"[test.py] Result for the experiment {i}: ")
            print(f"    Baseline throughput: {self.baseline_throughputs[i]}")
            print(f"    Groundtruth: {self.groundtruths[i]}")
            print(f"    Slowdown:    {self.slowdowns[i]}")
            print(f"    Predicted:   {self.predicteds[i]}")
            print(f"    Error Perc:  {self.errs[i]}")
    
    def __str__(self):
        max_key_length = max(len(key) for key in self.__dict__)
        return '\n'.join(f"{key.ljust(max_key_length)} : {value}" for key, value in self.__dict__.items())

def parse():
    parser = argparse.ArgumentParser()
    parser.add_argument("-b", "--benchmark",
                        help="the benchmark to run the experiment on",
                        type=str,
                        default="boutique")
    parser.add_argument("-t", "--num_threads",
                        help="the number of threads to run the experiment on",
                        type=int,
                        default=4)
    parser.add_argument("-c", "--num_conns", 
                        help="the number of connections to run the experiment on",
                        type=int,
                        default=128)
    # parser.add_argument("-d", "--duration",
    #                     help="the duration of the experiment",
    #                     type=int,
    #                     default=60)
    parser.add_argument("-x", "--target_service",
                        help="the target service to run the experiment on",
                        type=str,
                        default="cart")
    parser.add_argument("-r", "--request_type",
                        help="the request type to run the experiment on",
                        type=str,
                        default="mix")
    parser.add_argument("--repetitions",
                        help="the number of repetitions to run the experiment on",
                        type=int,
                        default=3)
    parser.add_argument("--num_exp", 
                        help="the number of experiments to run",
                        type=int,
                        default=10)
    parser.add_argument("--pre_run",
                        help="whether to run the experiment with prerun",
                        action="store_true")
    parser.add_argument("--range", nargs=2, type=int, default=[0, 1000])
    parser.add_argument("--num_req", type=int, default=50000)
    parser.add_argument("--clien_cpu_quota", type=int, default=-1)
    parser.add_argument("--random_seed", type=int, default=1234)
    parser.add_argument("--poker_batch", type=int, default=20000000)
    parser.add_argument("--poker_batch_req", type=int, default=100)
    parser.add_argument("--poker_relative_batch", type=int, default=20000000)
    args = parser.parse_args()
    return args

def main(args):
    runner = Runner(args)
    print(runner)
    runner.start()
    runner.print_summary()

if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    args = parse()
    main(args)
