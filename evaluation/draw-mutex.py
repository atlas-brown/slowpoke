import os
import numpy as np
import matplotlib.pyplot as plt
from pathlib import Path
from math import sqrt
import sys, re, ast, requests
import warnings

warnings.filterwarnings("ignore")

# # Set global font to Times New Roman
# plt.rcParams['font.family'] = 'Times New Roman'
# plt.rcParams['mathtext.fontset'] = 'custom'  # For math text
# plt.rcParams['mathtext.rm'] = 'Times New Roman'
# plt.rcParams['mathtext.it'] = 'Times New Roman:italic'
# plt.rcParams['mathtext.bf'] = 'Times New Roman:bold'

basedir = Path(os.path.dirname(os.path.abspath(__file__))).parent
imagedir = basedir / "plot"

def m_err(l):
    mean = sum(l)/len(l)
    err = sqrt(sum([(x-mean)*(x-mean) for x in l])/(len(l)-1))
    return mean, err

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

def load_data(results_dir):
    mutex_log = results_dir / "mutex-microbenchmark-lock.log"
    non_mutex_log = results_dir / "mutex-microbenchmark-no-lock.log"
    mutex_data = parse_results(mutex_log)
    non_mutex_data = parse_results(non_mutex_log)

    nomutex_truth = [non_mutex_data['groundtruth']]
    nomutex_predicted = [non_mutex_data['predicted']]
    mutex_truth = [mutex_data['groundtruth']]
    mutex_predicted = [mutex_data['predicted']]
    everything = [nomutex_truth, nomutex_predicted, mutex_truth, mutex_predicted]
    everything = [np.array(x) for x in everything]
    everymean = [x.T.mean(1) for x in everything]
    everyerr = [np.std(x.T, -1, ddof=1) for x in everything]
    return everymean, everyerr

if __name__ == '__main__':
    x = 800 - np.array(list(range(400, 800, 40)))
    means, errs = load_data(Path(sys.argv[1]))
    labels = ['No Mutex Groundtruth', 'No Mutex Predicted',
              'Mutex Groundtruth', 'Mutex Predicted']
    # plt.figure(figsize=(6, 2))
    marker = ['-o', '-o', '-s', '-s']
    colors = ["#a6611a", "#dfc27d", "#80cdc1", "#018571"]
    
    # Plot all the data first and store the line objects
    lines = []
    for k in [0, 1, 2, 3]:
        y, yerr = means[k], errs[k]
        line = plt.errorbar(x, y, yerr=yerr, fmt=marker[k], capsize=5, 
                           label=labels[k], color=colors[k])
        lines.append(line)
    
    plt.xlabel('Optimized Processing Time (us)', fontsize=12)
    plt.ylabel('Throughput (req/s)', fontsize=12)
    plt.xticks(x, fontsize=12)
    plt.yticks(fontsize=12, rotation=45)
    plt.ylim(0, 3500)
    
    # Create first legend (upper left) - first two items
    legend1 = plt.legend([lines[0], lines[1]], labels[:2], 
                        loc='upper left', fontsize=12)
    
    # Add the first legend manually to the current Axes
    plt.gca().add_artist(legend1)
    
    # Create second legend (bottom right) - last two items
    plt.legend([lines[2], lines[3]], labels[2:], 
               loc='lower right', fontsize=12)

    plt.savefig(f'/var/www/html/mutex_effectiveness.pdf')
    print(f'Result for mutex_effectiveness is available at')
    print(f'http://{IP}/mutex_effectiveness.pdf')
    plt.close()
