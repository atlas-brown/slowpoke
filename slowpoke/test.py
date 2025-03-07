#!python3

import os
import subprocess
import copy
import argparse
import config

REQ_COUNT = 12500

class Runner:
    def __init__(self, args):
        self.benchmark = args.benchmark
        self.num_threads = args.num_threads
        self.num_conns = args.num_conns
        self.target_service = args.target_service
        self.request_type = args.request_type
        self.request_ratio = config.get_request_ratio(self.benchmark, self.request_type)
        self.baseline_service_processing_time = config.get_baseline_service_processing_time(self.benchmark)
        self.repetitions = args.repetitions
        self.target_processing_time_range = args.range
        self.target_num_exp = args.num_exp
        self.duration = args.duration
        self.pre_run = args.pre_run
        self.groundtruths = []
        self.slowdowns = []
        self.predicteds = []
        self.errs = []
    
    def exp(self, service_delay):
        env = os.environ.copy()
        for service, delay in service_delay.items():
            env[f"SLOWPOKE_DELAY_MICROS_{service.upper()}"] = str(delay)
        env["SLOWPOKE_PRERUN"] = str(self.pre_run) # Disable request counting during normal execution
        cmd = f"bash run.sh {self.benchmark} {self.request_type} {self.num_threads} {self.num_conns} {self.duration}"
        print(f"[test.py] Running (pre_run: {self.pre_run}) workload {self.benchmark}/{self.request_type} request with the following configuration: {service_delay}", flush=True)
        process = subprocess.Popen(cmd, shell=True, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        print(f"[test.py] Executing cmd `{cmd}`:")
        throughput = -1
        times = []
        for line in process.stdout:
            line_output = line.decode().strip()
            print(f"    {line_output}", flush=True)
            time_prefix = 'stop time: '
            if line_output.startswith(time_prefix):
                times.append(float(line_output[len(time_prefix):]))
        throughput = self.num_threads * REQ_COUNT / (sum(times)/len(times))
        if process.wait() != 0:
            print(f"Error running {cmd}")
            for line in process.stderr:
                print(f"    {line.decode().strip()}", flush=True)
            raise Exception(f"Error running {cmd}")
        return throughput

    def run(self):
        service_delay = copy.deepcopy(self.baseline_service_processing_time)
        processing_time_diff = self.target_processing_time_range[1]-self.target_processing_time_range[0]
        processing_time_range = range(self.target_processing_time_range[0], self.target_processing_time_range[1], processing_time_diff//self.target_num_exp)

        # groundtruth
        groundtruth = []
        for p_t in processing_time_range:
            service_delay[self.target_service] = p_t
            res = self.exp(service_delay)
            while int(res) == 0:
                print("[test.py] Found 0 throughtput, rerun experiment")
                res = self.exp(service_delay)
            groundtruth.append(res)
        
        # slowdown
        slowdown = []
        predicted = []
        for p_t in processing_time_range:
            for service in service_delay:
                if service == self.target_service:
                    service_delay[service] = self.target_processing_time_range[1]
                else:
                    if self.request_ratio[service] == 0:
                        delay = 0
                    else:
                        delay = int(
                            ((self.target_processing_time_range[1] - p_t)*self.request_ratio[self.target_service]) / self.request_ratio[service]
                        )
                    service_delay[service] = self.baseline_service_processing_time[service] + delay
            res = self.exp(service_delay)
            while int(res) == 0:
                print("[test.py] Found 0 throughput, rerun experiment")
                res = self.exp(service_delay)
            slowdown.append(res)

            try:
                delay = (self.target_processing_time_range[1] - p_t)*self.request_ratio[self.target_service]
                predicted_throughput = 1000000/(1000000/slowdown[-1]-delay)
            except:
                print("[test.py] Error: Division by zero")
                predicted_throughput = -1
            predicted.append(predicted_throughput)
        
        err = [(predicted[i]-groundtruth[i])*100/groundtruth[i] for i in range(len(predicted))]
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
            print(f"    Groundtruth: {self.groundtruths[i]}")
            print(f"    Slowdown:    {self.slowdowns[i]}")
            print(f"    Predicted:   {self.predicteds[i]}")
            print(f"    Error Perc:  {self.errs[i]}")

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
    parser.add_argument("-d", "--duration",
                        help="the duration of the experiment",
                        type=int,
                        default=60)
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
    args = parser.parse_args()
    return args

def main(args):
    runner = Runner(args)
    print(f"[test.py] BENCHMARK:                        {runner.benchmark}")
    print(f"[test.py] REQUEST_TYPE:                     {runner.request_type}")
    print(f"[test.py] TARGET_SERVICE:                   {runner.target_service}")
    print(f"[test.py] NUM_THREADS:                      {runner.num_threads}")
    print(f"[test.py] NUM_CONN:                         {runner.num_conns}")
    print(f"[test.py] DURATION:                         {runner.duration}")
    print(f"[test.py] REPETITIONS:                      {runner.repetitions}")
    print(f"[test.py] REQUEST_RATIO:                    {runner.request_ratio}")
    print(f"[test.py] BASELINE_PROCESSING_TIME:         {runner.baseline_service_processing_time}")
    print(f"[test.py] TARGET_PROCESSING_TIME_RANGE:     {runner.target_processing_time_range}")
    print(f"[test.py] TARGET_NUM_EXP:                   {runner.target_num_exp}")
    runner.start()
    runner.print_summary()

if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    args = parse()
    main(args)
