 sudo cpupower frequency-set -g performance
  echo 1 |sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo
