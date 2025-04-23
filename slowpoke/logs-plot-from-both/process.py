import os
import re
import matplotlib.pyplot as plt
import numpy as np
from matplotlib.colors import TABLEAU_COLORS

def parse_log_file(filepath):
    """Parse a single log file and return sleep intervals"""
    service_name = os.path.basename(filepath).split('-')[0]
    intervals = []
    intervals_pc = []
    with open(filepath, 'r') as f:
        count = 0
        base_time = -1
        for line in f:
            if count > 100:
                break
            if "Received 0 ns" in line:
                continue
            match = re.search(
                r'started at (\d+), ended at (\d+)', 
                line.strip()
            )
            if match:
                if count == 0:
                    base_time = int(match.group(1))
                start = (int(match.group(1)) - base_time) / 100000000  # ns to μs
                end = start  + 0.1    # ns to μs
                intervals.append({
                    'service': service_name,
                    'start': start,
                    'end': end,
                    'duration': end - start
                })
                count += 1
        count = 0
        for line in f:
            if count > 100:
                break
            # if "Received 0 ns" in line:
                # continue
            match = re.search(
                r'started pc at (\d+), ended pc at (\d+)', 
                line.strip()
            )
            if match:
                if count == 0:
                    base_time = int(match.group(1))
                start = (int(match.group(1))-base_time) / 100000000  # ns to μs
                end = start  + 0.1    # ns to μs
                intervals_pc.append({
                    'service': service_name,
                    'start': start,
                    'end': end,
                    'duration': end - start
                })
                count += 1
        for i in range(10):
            print(f"intervals: {intervals[i]}")
            print(f"intervals_pc: {intervals_pc[i]}")
    return intervals, intervals_pc

def plot_timeline(intervals, intervals_pc):
    """Plot sleep periods as horizontal line segments"""
    if not intervals:
        print("No intervals found to plot")
        return
    
    services = sorted(list(set(i['service'] for i in intervals)))
    colors = list(TABLEAU_COLORS.values())
    
    fig, ax = plt.subplots(figsize=(12, 6))
    
    # Plot each service's sleep periods
    for y, service in enumerate(services, 1):
        service_intervals = [i for i in intervals if i['service'] == service]
        service_pc_intervals = [i for i in intervals_pc if i['service'] == service]
        color = colors[y % len(colors)]
        
        # Plot each sleep interval as a separate line segment
        for interval in service_intervals:
            ax.plot(
                [interval['start'], interval['end']],
                [y+0.1, y+0.1],
                color=color,
                linewidth=3,
                solid_capstyle='butt',
                label=service if interval == service_intervals[0] else None
            )
        for interval in service_pc_intervals:
            ax.plot(
                [interval['start'], interval['end']],
                [y-0.1, y-0.1],
                color="black",
                linewidth=3,
                linestyle='dashed',
                solid_capstyle='butt',
                label=service if interval == service_pc_intervals[0] else None
            )
    
    # Formatting
    ax.set_yticks(range(1, len(services)+1))
    ax.set_yticklabels(services)
    ax.set_xlabel('Time (μs)')
    ax.set_title('Service Sleep Periods')
    
    # Auto-scale time axis
    min_time = min(i['start'] for i in intervals)
    max_time = max(i['end'] for i in intervals)
    time_range = max_time - min_time
    
    # if time_range > 1e6:  # > 1 second
    #     ax.set_xlabel('Time (seconds)')
    #     ax.xaxis.set_major_formatter(plt.FuncFormatter(lambda x, _: f'{x/1e6:.3f}'))
    # el
    if time_range > 1e3:  # > 1 millisecond
        ax.set_xlabel('Time (ms)')
        ax.xaxis.set_major_formatter(plt.FuncFormatter(lambda x, _: f'{x/1e3:.3f}'))
    
    # Clean up borders
    for spine in ['top', 'right', 'left']:
        ax.spines[spine].set_visible(False)
    
    ax.grid(True, axis='x')
    ax.legend(loc='upper right')
    
    plt.tight_layout()
    plt.savefig('service_sleep_periods.png', dpi=300, bbox_inches='tight')
    plt.close()

def main(log_dir):
    all_intervals = []
    all_intervals_pc = []
    
    for filename in os.listdir(log_dir):
        if filename.startswith('service') and filename.endswith('.log'):
            filepath = os.path.join(log_dir, filename)
            intervals, intervals_pc = parse_log_file(filepath)
            all_intervals.extend(intervals)
            all_intervals_pc.extend(intervals_pc)
    if all_intervals:
        all_intervals.sort(key=lambda x: x['start'])
        all_intervals_pc.sort(key=lambda x: x['start'])
        plot_timeline(all_intervals, all_intervals_pc)
    else:
        print("No valid intervals found")

if __name__ == '__main__':
    import sys
    log_dir = sys.argv[1] if len(sys.argv) > 1 else "."
    main(log_dir)
