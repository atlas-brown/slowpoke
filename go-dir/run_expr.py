import os
import json
import time
from dataclasses import dataclass
from pathlib import Path
from subprocess import Popen, run, PIPE

basedir = Path(os.path.dirname(os.path.abspath(__file__)))
datadir = basedir / "data"
yamldir = basedir / "yamls"

PODNAME = 'testtime'

@dataclass
class Expr:
    kind: str
    num: int
    lim: bool

    def name(self):
        l = "lim" if self.lim else "nolim"
        return f'{self.kind}-{self.num}-{l}'

    def yaml_file(self):
        return f'{self.name()}.yaml'
    
def build_image():
    result = run("docker build . -t sleepymug/golang-timing-test --push".split())
    return result.returncode == 0

def create_yaml(expr):
    template_path = yamldir / "template.yaml"
    template = open(template_path).read()
    with open(yamldir / f"{expr.yaml_file()}", 'w') as f:
        f.write(template.replace('{{kind}}', expr.kind)
                        .replace('{{num}}', str(expr.num))
                        .replace('{{cpu}}', 'cpu: 1000m' if expr.lim else ''))

def runcmd_capture(cmd):
    print(f'[?] {cmd}')
    result = run(cmd, shell=True, stdout=PIPE, text=True)
    return result.stdout

def runcmd(cmd):
    print(f'[+] {cmd}')
    run(cmd, shell=True)
        
def remove_pod():
    runcmd(f'kubectl delete pod {PODNAME}')
    
def wait_pod():
    while True:
        out = runcmd_capture(f'kubectl get pod {PODNAME}'
               ' -o jsonpath="{.status.containerStatuses[0].state.terminated.reason}"')
        if out == 'Completed':
            break
        time.sleep(2)
    
def run_pod(yaml_file):
    runcmd(f'kubectl apply -f {yaml_file}')

def read_pod_log():
    output = runcmd_capture(f'kubectl logs {PODNAME}')
    print(output)
    return int(output.strip())
    
def run_expr(expr):
    create_yaml(expr)
    # remove_pod()
    run_pod(yamldir / expr.yaml_file())
    wait_pod()
    ret = read_pod_log()
    remove_pod()
    return ret
    
def sanitize_env():
    run(f'kubectl delete pod {PODNAME}', shell=True, stdout=PIPE, stderr=PIPE, text=True)
    datadir.mkdir(parents=True, exist_ok=True)
    
if __name__ == '__main__':
    os.chdir(basedir)
    # r = build_image()
    # if not r:
    #     print('docker image build problem')
    #     exit(1)
    sanitize_env()
    for kind in ['getpid', 'spin', 'spinx']:
        for num in range(11):
            for lim in [True, False]:
                results = []
                expr = Expr(kind, num*1000, lim)
                for repeat in range(10):
                    res = run_expr(expr)
                    results.append(str(res))
                with open(datadir / expr.name(), 'w') as f:
                    f.write('\n'.join(results))
                    f.write('\n')
