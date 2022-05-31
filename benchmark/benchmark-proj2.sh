#!/bin/bash
#
#SBATCH --mail-user=wjsjr@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj2_benchmark
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/wjsjr/Parallel/proj2-wjs97/proj2/benchmark
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=180:00


module load golang/1.16.2
python3 benchmark.py
