import os
import pandas as pd

# Set display options to avoid truncation
pd.set_option('display.max_rows', None)  # or use a high number instead of None
pd.set_option('display.max_columns', None)  # adjust according to your needs
pd.set_option('display.width', None)  # adjust width for your display
pd.set_option('display.max_colwidth', None)  # or a high number to avoid truncation of column content

directories = [d for d in os.listdir('./results') if os.path.isdir(f'./results/{d}')]
latest_directory = sorted(directories)[-1]

csv_files = [f for f in os.listdir(f'./results/{latest_directory}') if f.endswith('.csv')]

def sort_key(f):
    return int(f.split('_')[1])

csv_files.sort(key=sort_key)

for file in csv_files:
    df = pd.read_csv(f'./results/{latest_directory}/{file}')
    df_without_best = df.drop(columns=['BestSolution'])
    print(f'File: {file}')  # Adjust slicing as needed to display the desired part of the filename
    print(df_without_best)
    print()
    best_solution = [int(x) for x in df.loc[df['BestSolution'].idxmin(), "BestSolution"].strip('[]').split()]
    print("Best solution")
    print(best_solution)
    print()
