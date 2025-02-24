import subprocess

process = subprocess.run(
    ["bash", "-c", 'kubectl get pods | cut -f 1 -d " " | grep -vE "ubuntu|NAME|redis"'],
    capture_output=True,
    text=True,
    check=True,
)
stdout = process.stdout.strip()
counters = {}
for pod in stdout.splitlines():
    log_process = subprocess.run(
        ["bash", "-c", f'kubectl logs {pod}'],
        capture_output=True,
        text=True,
        check=True,
    )
    logs_stdout = log_process.stdout.strip()
    service = pod.split("-")[0]
    counters[service] = 0
    for line in logs_stdout.splitlines():
        if "[Total]" in line:
            value = int(line.split("[Total]")[1].strip())
            counters[service] += value
ratios = {
    service: counters[service] / counters["frontend"] for service in counters
}
print(counters)
print(ratios)