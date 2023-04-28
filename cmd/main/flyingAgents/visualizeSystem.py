import matplotlib.pyplot as plt
import matplotlib.animation as animation
import matplotlib.patches as patches
import pandas as pd
import numpy as np
from GonumMatrixIO import GonumIO
import argparse

simulationData = pd.read_parquet("data/BestAgentSimulation.pq")
animationSavePath = "data/animation.mp4"

GAME_DIMENSION = 100
AGENT_RADIUS = 5
TARGET_LOCATION_RADIUS = 2

parser = argparse.ArgumentParser(description="Flying Agents Argument Parser")
parser.add_argument("--save", help="Save animation to file, rather than showing", action="store_true")
parser.add_argument("--numFrames", help="Determine the number of frames to render. If not given, render the entire simulation", action="store", type=int, default=None)
args = parser.parse_args()

if args.numFrames == None or args.numFrames > len(simulationData):
    numFrames = len(simulationData)
else:
    numFrames = args.numFrames

print(f"BEST CHROMOSOME:\n{GonumIO.loadMatrix('data/bestAgentChromosome.bin')}")
    

fig, ax = plt.subplots(figsize=(10,10))
plt.title("Genetic Algorithm: Flying Agents Visualization")
plt.xlim(-GAME_DIMENSION, GAME_DIMENSION)
plt.ylim(-GAME_DIMENSION, GAME_DIMENSION)
agent = patches.Circle((0,0), radius=AGENT_RADIUS, color='k')
targetLocation = patches.Circle((0,0), radius=TARGET_LOCATION_RADIUS, color='r')
ax.add_patch(agent)
ax.add_patch(targetLocation)

def update(index):

    stateVector = simulationData.loc[index, "StateVector"]
    agentX = stateVector[0]
    agentY = stateVector[1]
    agentVelX = stateVector[2]
    agentVelY = stateVector[3]
    targetLocationX = stateVector[4]
    targetLocationY = stateVector[5]

    print(agentX, agentY, agentVelX, agentVelY, targetLocationX, targetLocationY)
    agent.set_center((agentX, agentY))
    targetLocation.set_center((targetLocationX, targetLocationY))
    return agent, targetLocation,

anim = animation.FuncAnimation(fig, update, frames=numFrames, interval=1)
if args.save:
    writer = animation.FFMpegWriter(fps=60)
    anim.save(animationSavePath, writer=writer)
else:
    plt.show()