import os
import numpy as np

os.chdir(os.path.dirname(os.path.abspath(__file__)))
OUTPUT_FILE="statistics.txt"

def parse_result_file(filename):
    groundtruths = [[] for _ in range(10)]
    predictions = [[] for _ in range(10)]
    baselines = []
    errors = []
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
    # sort each entry, pick the middle three
    for i in range(len(groundtruths)):
        groundtruths[i].sort()
        predictions[i].sort()
        start_idx = len(groundtruths[i]) // 2 - 1
        groundtruths[i] = groundtruths[i][start_idx:start_idx+3]
        predictions[i] = predictions[i][start_idx:start_idx+3]
        # calculate error
        for j in range(len(groundtruths[i])):
            errors.append(100 * (predictions[i][j] - groundtruths[i][j]) / groundtruths[i][j])
    baselines.sort()
    start_idx = len(baselines) // 2 - 1
    baselines = baselines[start_idx:start_idx+3]
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

    print(f"File: {filename}")
    print(f"    Mean Error % (Raw): {mean_error:.2f}")
    print(f"    Min Error % (Raw): {min_error:.2f}")
    print(f"    Max Error % (Raw) : {max_error:.2f}")
    print(f"    Std (Raw): {np.std(errors):.2f}")
    print(f"    Median Error % (Raw): {np.median(errors):.2f}")
    print(f"    Mean Error % (Abs): {mean_error_abs:.2f}")
    print(f"    Min Error % (Abs): {min_error_abs:.2f}")
    print(f"    Max Error % (Abs) : {max_error_abs:.2f}")
    print(f"    Std (Abs): {np.std(errors_abs):.2f}")
    # also print the ranges of the groundtruth and predicted values
    flattened_groundtruths = np.array(groundtruths).ravel()
    flattened_predictions = np.array(predictions).ravel()
    print(f"    Groundtruth Range (Rps): {np.min(flattened_groundtruths)} - {np.max(flattened_groundtruths)}")
    print(f"    Predicted Range (Rps): {np.min(flattened_predictions)} - {np.max(flattened_predictions)}")

    return experiments, errors

# Define the processing time range
opt_range = np.arange(0, 100, 10)
all_errs = []

# Process each log file
for file_name in os.listdir("all-results"):
    experiments, err = parse_result_file(f"./all-results/{file_name}")
    all_errs.append(err.flatten())

errors = np.array(all_errs)
errors_abs = np.abs(errors)

mean_error = np.mean(errors)
min_error = np.min(errors)
max_error = np.max(errors)
mean_error_abs = np.mean(errors_abs)
min_error_abs = np.min(errors_abs)
max_error_abs = np.max(errors_abs)
print("Overall Statistics:")
print(f"    Mean Error % (Raw): {mean_error:.2f}")
print(f"    Min Error % (Raw): {min_error:.2f}")
print(f"    Max Error % (Raw) : {max_error:.2f}")
print(f"    Std (Raw): {np.std(errors):.2f}")
print(f"    Median Error % (Raw): {np.median(errors):.2f}")
print(f"    Mean Error % (Abs): {mean_error_abs:.2f}")
print(f"    Min Error % (Abs): {min_error_abs:.2f}")
print(f"    Max Error % (Abs) : {max_error_abs:.2f}")
print(f"    Std (Abs): {np.std(errors_abs):.2f}")