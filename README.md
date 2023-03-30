# Foosball Genetic Learning

This project is an attempt to build a genetic learning model around a simulated foosball table. We aim to train initially random agents into component foosball players using nothing more than repeated matches against one another. 

Right now we are in the very early planning stages - ensuring everyone is on the same page before starting implementations.

Among the specifics yet to be determined:

- The exact implementation of the system/environment
    - Should `AgentAction` determine foosball rod position, velocity, or acceleration (tradeoff between simulation complexity and accuracy)
    - What should the scoring function of the Agent look like? (number of goals, time ball spent in opposition half, amount of work moving rods...)
- The exact implementation of the genetic algorithm
    - Parent selection method (n-Best Agents, Score Distribution, ...)
    - Chromosome combination method (n-point crossover, other?)
    - Mutation rate(s) 
    - Fittest Carryover