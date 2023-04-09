import matplotlib.pyplot as plt
import pandas as pd
import numpy as np

GAME_X_DIMENSION = 1.0
GAME_Y_DIMENSION = 0.5
PADDLE_SIZE = 0.2

simulationData = pd.read_parquet("data/BestAgentSimulation.pq")

def parseStateVector(stateVector):
    ballX = stateVector[0]
    ballY = stateVector[1]
    ballXVelocity = stateVector[2]
    ballYVelocity = stateVector[3]
    paddle0Position = stateVector[4]
    paddle1Position = stateVector[5]
    # print(ballX,ballY,ballXVelocity,ballYVelocity,paddle0Position,paddle1Position,)

    ball.set_xdata([ballX])
    ball.set_ydata([ballY])
    paddle0.set_ydata([paddle0Position-PADDLE_SIZE, paddle0Position+PADDLE_SIZE])
    paddle1.set_ydata([paddle1Position-PADDLE_SIZE, paddle1Position+PADDLE_SIZE])
    plt.draw()
    plt.pause(0.01)

fig, ax = plt.subplots()
plt.xlim(-1.1*GAME_X_DIMENSION, 1.1*GAME_X_DIMENSION)
plt.ylim(-1.1*GAME_Y_DIMENSION, 1.1*GAME_Y_DIMENSION)
paddle0, = ax.plot([-GAME_X_DIMENSION,-GAME_X_DIMENSION],[-PADDLE_SIZE, PADDLE_SIZE])
paddle1, = ax.plot([GAME_X_DIMENSION,GAME_X_DIMENSION],[-PADDLE_SIZE, PADDLE_SIZE])
ball, = ax.plot(0,0,marker="o", markersize=10)

for rowindex, row in simulationData.iterrows():
    parseStateVector(row["StateVector"])