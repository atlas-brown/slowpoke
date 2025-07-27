#!/usr/bin/env python3
import os
import matplotlib.pyplot as plt
import numpy as np
from scipy.interpolate import make_interp_spline
import argparse
import requests
import warnings
from pathlib import Path
import sys
warnings.filterwarnings("ignore", category=RuntimeWarning)

def get_public_ip():
    resp = requests.get('http://checkip.amazonaws.com')
    resp.raise_for_status()            # optional: will raise an error if the request failed
    return resp.text.strip()          # strip() removes the trailing newline

IP = get_public_ip()

# Set global font to Times New Roman
# plt.rcParams['font.family'] = 'Times New Roman'
# plt.rcParams['mathtext.fontset'] = 'custom'  # For math text
# plt.rcParams['mathtext.rm'] = 'Times New Roman'
# plt.rcParams['mathtext.it'] = 'Times New Roman:italic'
# plt.rcParams['mathtext.bf'] = 'Times New Roman:bold'

parser = argparse.ArgumentParser()
parser.add_argument('-n', type=int, default=10, help='Number of optimizations')
args = parser.parse_args()
NUM_OPTIMIZATIONS = args.n

def filter_outliers(baselines):
    # Sort the baselines
    baselines.sort()

    # Calculate median and MAD (Median Absolute Deviation)
    median = np.median(baselines)
    mad = np.median(np.abs(baselines - median))

    # Define a threshold (e.g., 3 MADs)
    threshold = 3


    # Filter out outliers
    filtered_baselines = [x for x in baselines if np.abs(x - median) / mad < threshold]

    # If fewer than 3 values, pick the closest remaining from original baselines
    if len(filtered_baselines) < 3:
        # Median of the filtered set (or original if empty)
        target_median = np.median(filtered_baselines) if filtered_baselines else median
        # Values not in filtered_baselines (i.e., outliers)
        remaining_values = [x for x in baselines if x not in filtered_baselines]
        # Compute distances to the target median
        distances = [abs(x - target_median) for x in remaining_values]
        # Get the closest missing values (up to 3 - len(filtered_baselines))
        closest_indices = np.argsort(distances)[:3 - len(filtered_baselines)]
        # Add them to the filtered list
        filtered_baselines += [remaining_values[i] for i in closest_indices]
        # Sort the final selection (optional)
        filtered_baselines.sort()

    # If 3 or more, pick the middle three
    if len(filtered_baselines) >= 3:
        start_idx = len(filtered_baselines) // 2 - 1
        filtered_baselines = filtered_baselines[start_idx:start_idx + 3]
    return filtered_baselines

def parse_result_file(filename):
    groundtruths = [[] for _ in range(NUM_OPTIMIZATIONS)]
    predictions = [[] for _ in range(NUM_OPTIMIZATIONS)]
    baselines = []
    errors = []
    diff_errors = []
    with open(filename, 'r') as file:
        lines = file.readlines()
        i = 0
        while i < len(lines):
            if "Result for the experiment" in lines[i]:
                baseline = float(lines[i+1].split(":")[1].strip())
                groundtruth = list(map(float, lines[i+2].split(":")[1].strip().strip('[]').split(',')))[::-1]
                predicted = list(map(float, lines[i+4].split(":")[1].strip().strip('[]').split(',')))[::-1]
                # error = list(map(float, lines[i+5].split(":")[1].strip().strip('[]').split(',')))[::-1]
                # errors.extend(error)
                # experiments.append((baseline, groundtruth, predicted))
                for j in range(len(groundtruth)):
                    groundtruths[j].append(groundtruth[j])
                    predictions[j].append(predicted[j])
                baselines.append(baseline)
                i += 6
            else:
                i += 1
    groundtruth_variance = []
    groundtruth_stdev = []
    groundtruth_diff = []
    groundtruth_mean = []
    # sort each entry, pick the middle three
    for i in range(len(groundtruths)):
        groundtruths[i].sort()
        predictions[i].sort()
        start_idx = len(groundtruths[i]) // 2 - 1
        groundtruths[i] = groundtruths[i][start_idx:start_idx+3]
        predictions[i] = predictions[i][start_idx:start_idx+3]
        variance = np.var(groundtruths[i])
        groundtruth_variance.append(variance)
        stdev = np.std(groundtruths[i])
        groundtruth_stdev.append(stdev)
        diff = (np.max(groundtruths[i]) - np.min(groundtruths[i]))
        groundtruth_diff.append(diff)
        groundtruth_mean.append(np.mean(groundtruths[i]))
        # groundtruths[i] = filter_outliers(groundtruths[i])
        # predictions[i] = filter_outliers(predictions[i])
        # calculate error
        for j in range(len(groundtruths[i])):
            errors.append(100 * (predictions[i][j] - groundtruths[i][j]) / groundtruths[i][j])
            diff_errors.append((predictions[i][j] - groundtruths[i][j]))
    baselines.sort()
    # start_idx = len(baselines) // 2 - 1
    # baselines = baselines[start_idx:start_idx+3]
    baselines = filter_outliers(baselines)

    # print(len(baselines), len(groundtruths), len(predictions))
    # construct experiments
    experiments = []
    for i in range(len(baselines)):
        experiments.append((baselines[i], [groundtruths[j][i] for j in range(len(groundtruths))], [predictions[j][i] for j in range(len(predictions))]))
    # calculate average error, min and max, using absolute error
    errors = np.array(errors)
    errors_abs = np.abs(errors)

    mean_error = np.mean(errors)
    min_error = np.min(errors)
    max_error = np.max(errors)
    mean_error_abs = np.mean(errors_abs)
    min_error_abs = np.min(errors_abs)
    max_error_abs = np.max(errors_abs)

    diff_errors = np.array(diff_errors)
    mean_diff_error = np.mean(diff_errors)
    stdev_diff_error = np.std(diff_errors)
    rsme = np.sqrt(np.mean(errors**2))

    # print(f"File: {filename}")
    # print(f"    Mean Error % (Raw): {mean_error:.2f}")
    # print(f"    Min Error % (Raw): {min_error:.2f}")
    # print(f"    Max Error % (Raw) : {max_error:.2f}")
    # print(f"    Std (Raw): {np.std(errors):.2f}")
    # print(f"    Median Error % (Raw): {np.median(errors):.2f}")
    # print(f"    RSME % (Raw): {rsme:.2f}")
    # print(f"    Mean Error % (Abs): {mean_error_abs:.2f}")
    # print(f"    Min Error % (Abs): {min_error_abs:.2f}")
    # print(f"    Max Error % (Abs) : {max_error_abs:.2f}")
    # print(f"    Std (Abs): {np.std(errors_abs):.2f}")
    # print(f"    Mean Diff Error % (Raw): {mean_diff_error:.2f}")
    # print(f"    Std Diff Error % (Raw): {stdev_diff_error:.2f}")
    # print(f"    Min Diff Error % (Raw): {np.min(diff_errors):.2f}")
    # print(f"    Max Diff Error % (Raw): {np.max(diff_errors):.2f}")
    # # also print the ranges of the groundtruth and predicted values
    # flattened_groundtruths = np.array(groundtruths).ravel()
    # flattened_predictions = np.array(predictions).ravel()
    # print(f"    Groundtruth Range (Rps): {np.min(flattened_groundtruths)} - {np.max(flattened_groundtruths)}")
    # print(f"    Predicted Range (Rps): {np.min(flattened_predictions)} - {np.max(flattened_predictions)}")
    # print(f"    Groundtruth Variance: {groundtruth_variance}")
    # print(f"    Groundtruth Stdev: {groundtruth_stdev}")
    # print(f"    Groundtruth Diff: {groundtruth_diff}")
    # print(f"    Groundtruth Mean: {groundtruth_mean}")
    # print(f"    Groundtruth Diff range: {np.min(groundtruth_diff)} - {np.max(groundtruth_diff)}")

    return experiments, errors

result_dir = Path(sys.argv[1])

# List of log files to process
log_files = [
    result_dir / "hotel_medium.log",
    result_dir / "movie_medium.log",
    result_dir / "boutique_medium.log",
    result_dir / "social_medium.log"
]

benchmark_names = [
    "HotelReservation",
    "MovieReview",
    "OnlineBoutique",
    "SocialMedia"
]

# Create 1x4 subplots
fig, axs = plt.subplots(1, 4, figsize=(12, 2.5))  # Slightly wider to accommodate error bars
fig.tight_layout(pad=3.0)

# Define the processing time range
# opt_range = np.arange(10, 101, 10)
all_errs = []
target_ratios = [0.8, 0.9050900514, 0.45175405908969923, 0.7005428505]

# Process each log file
for idx, (file_name, bench_name) in enumerate(zip(log_files, benchmark_names)):
    experiments, err = parse_result_file(file_name)

    # grep baselines
    baselines = [exp[0] for exp in experiments]
    avg_baseline = np.mean(baselines)
    base_processing_time = 2e6 * target_ratios[idx] / avg_baseline + 1000
    step = int(1000/NUM_OPTIMIZATIONS)
    opt_range = np.array([(i) * 100 / base_processing_time for i in range(step, 1001, step)])

    all_errs.append(err.flatten())
    ax = axs[idx]
    
    # Prepare data for error bars
    groundtruth_values = np.array([exp[1] for exp in experiments])
    predicted_values = np.array([exp[2] for exp in experiments])
    
    # Calculate means and standard deviations
    gt_mean = np.mean(groundtruth_values, axis=0)
    gt_std = np.std(groundtruth_values, axis=0)
    pred_mean = np.mean(predicted_values, axis=0)
    pred_std = np.std(predicted_values, axis=0)

    y_data = np.concatenate([groundtruth_values.flatten(), predicted_values.flatten()])
    median_y = np.median(y_data)
    upper_lim = 1.05 * np.max(y_data)
    lower_lim = 0.95 *np.min(y_data)

    # Plot with error bars
    ax.errorbar(opt_range, gt_mean, yerr=[gt_mean-np.min(groundtruth_values, axis=0), np.max(groundtruth_values, axis=0)-gt_mean], 
                fmt='o-', color='black', markersize=4, capsize=2, linewidth=1,
                label='Groundtruth' if idx == 0 else "")
    
    ax.errorbar(opt_range, pred_mean, yerr=[pred_mean-np.min(predicted_values, axis=0), np.max(predicted_values, axis=0)-pred_mean],
                fmt='s--', color='#1f77b4', markersize=4, capsize=2, linewidth=1,
                label='Predicted' if idx == 0 else "")

    # Set dynamic y-axis limits
    # print(lower_lim, upper_lim)
    ax.set_ylim(lower_lim, upper_lim)
    
    # Grid and labels
    ax.grid(True, linestyle=':', alpha=0.6)
    ax.set_title(benchmark_names[idx], fontsize=12)
    
    # Y-axis (only numbers, no label)
    ax.set_ylabel('End Throughput (req/s)' if idx == 0 else '', fontsize=12)
    ax.set_yticks(np.linspace(ax.get_yticks()[0], ax.get_yticks()[-1], num=3))
    
    # Rotate y-axis labels
    for label in ax.get_yticklabels():
        label.set_rotation(45)
        label.set_ha('right')
        label.set_rotation_mode('anchor')

errors = np.array(all_errs)
rsme = np.sqrt(np.mean(errors**2))
errors_abs = np.abs(errors)

mean_error = np.mean(errors)
min_error = np.min(errors)
max_error = np.max(errors)
mean_error_abs = np.mean(errors_abs)
min_error_abs = np.min(errors_abs)
max_error_abs = np.max(errors_abs)
# print("Overall Statistics:")
# print(f"    Mean Error % (Raw): {mean_error:.2f}")
# print(f"    Min Error % (Raw): {min_error:.2f}")
# print(f"    Max Error % (Raw) : {max_error:.2f}")
# print(f"    Std (Raw): {np.std(errors):.2f}")
# print(f"    RSME % (Raw): {rsme:.2f}")
# print(f"    Mean Error % (Abs): {mean_error_abs:.2f}")
# print(f"    Mean Error % (Abs): {mean_error_abs:.2f}")
# print(f"    Min Error % (Abs): {min_error_abs:.2f}")
# print(f"    Max Error % (Abs) : {max_error_abs:.2f}")
# print(f"    Std (Abs): {np.std(errors_abs):.2f}")


# Add centered x-axis label
# fig.text(0.5, 0.02, 'Optimization Percentage (%)', ha='center', va='center', fontsize=10)

# Create custom legend in one row
legend_elements = [
    plt.Rectangle((0,0), 1, 1, fc='white', alpha=0.3, label='X-axis: Optimization Percentage (%)'),
    plt.Line2D([0], [0], marker='o', color='black', label='Groundtruth',
              markersize=6, linestyle='-', linewidth=1),
    plt.Line2D([0], [0], marker='s', color='#1f77b4', label='Predicted',
              markersize=6, linestyle='--', linewidth=1),
]

# Position legend below the plots (single row)
fig.legend(
    handles=legend_elements,
    loc='lower center',
    ncol=3,  # Ensures legend items stay in one row
    bbox_to_anchor=(0.5, -0.02),  # Adjust vertical position if needed
    frameon=False,
    fontsize=12
)

plt.subplots_adjust(bottom=0.25, top=0.85, wspace=0.18)  # Adjusted spacing
fig_file_name = os.path.basename(__file__).replace('.py', '.pdf')
os.makedirs("plot", exist_ok=True)
# plt.savefig(f"plot/{fig_file_name}", bbox_inches='tight', dpi=300)
plt.savefig(f'/var/www/html/{fig_file_name}', bbox_inches='tight', dpi=300)
print(f'Result for {fig_file_name} is available at')
print(f'http://{IP}/{fig_file_name}')
plt.close()
