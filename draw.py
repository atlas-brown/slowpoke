import sys
import requests
from pathlib import Path
import re, ast
import matplotlib.pyplot as plt

import requests

def get_public_ip():
    resp = requests.get('http://checkip.amazonaws.com')
    resp.raise_for_status()            # optional: will raise an error if the request failed
    return resp.text.strip()          # strip() removes the trailing newline

IP = get_public_ip()

def parse_results(filename):
    with open(filename, 'r') as f:
        text = f.read()

    # regex patterns for each field
    patterns = {
        'baseline':    r"Baseline throughput:\s*([0-9]+\.[0-9]+)",
        'groundtruth': r"Groundtruth:\s*(\[[^\]]+\])",
        'slowdown':    r"Slowdown:\s*(\[[^\]]+\])",
        'predicted':   r"Predicted:\s*(\[[^\]]+\])",
        'error_perc':  r"Error Perc:\s*(\[[^\]]+\])",
    }

    data = {}
    for key, pat in patterns.items():
        m = re.search(pat, text)
        if not m:
            raise ValueError(f"Could not find `{key}` in {filename}")
        s = m.group(1)
        # lists vs single floats
        if s.startswith('['):
            data[key] = ast.literal_eval(s)
        else:
            data[key] = float(s)
    return data

if __name__ == "__main__":
    result_dir = Path(sys.argv[1])
    for logfile in result_dir.glob("*.log"):
        data = parse_results(logfile)

        groundtruth = data['groundtruth']
        predicted   = data['predicted']
        slowdown    = data['slowdown']

        x = list(range(1, len(groundtruth) + 1))
        # Or to plot vs slowdown factors, uncomment:
        # x = slowdown
        x = [e * 20 for e in x]

        plt.figure()
        plt.plot(x, list(reversed(groundtruth)), label='Groundtruth', marker='o')
        plt.plot(x, list(reversed(predicted)),   label='Predicted',   marker='s')
        y_max = max([max(groundtruth), max(predicted)]) * 1.2
        plt.xlabel('Optimized processing time (%)')            # or 'Slowdown factor' if using slowdown
        plt.ylabel('Throughput')
        plt.ylim(0, y_max)
        plt.title('Predicted vs. Groundtruth Throughput')
        plt.legend()
        plt.tight_layout()
        figure_name = str(logfile.name).replace(".log", ".png")
        plt.savefig(f'/var/www/html/{figure_name}')
        print(f'Result for {logfile} is available at')
        print(f'http://{IP}/{figure_name}')
        plt.close()
