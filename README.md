# Foosball Genetic Learning

![Go Test](https://github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/actions/workflows/goTest.yaml/badge.svg)

This project is an attempt to build a genetic learning model around a simulated foosball table. We aim to train initially random agents into component foosball players using nothing more than repeated matches against one another. 

If you want to join this project, or just ask questions, feel free to email the members involved!

- Hayden McAlister: mcaha814@student.otago.ac.nz

#### Planning

Among the specifics yet to be determined:

---
On the Genetic Learning side of things:

- The exact implementation of the genetic algorithm
    - Parent selection method (n-Best Agents, Score Distribution, ...)
    - Chromosome combination method (n-point crossover, other?)
    - Mutation rate(s) 
    - Fittest Carryover

---
On the Foosball side of things:

- The exact implementation of the system/environment
    - Should `AgentAction` determine foosball rod position, velocity, or acceleration (tradeoff between simulation complexity and accuracy)
    - What should the scoring function of the Agent look like? (number of goals, time ball spent in opposition half, amount of work moving rods...)


#### Further Readings:

If you are interested in getting involved with this project (or want to learn more about genetic algorithms + learning), here are some helpful links:

- [Wikipedia: Genetic Algorithm](https://en.wikipedia.org/wiki/Genetic_algorithm)
- [TowardsDataScience Post: Introduction to Genetic Algorithms](https://towardsdatascience.com/introduction-to-genetic-algorithms-including-example-code-e396e98d8bf3)
- [Wikipedia: Crossover](https://en.wikipedia.org/wiki/Crossover_(genetic_algorithm))
- [Wikipedia: Mutation](https://en.wikipedia.org/wiki/Mutation_(genetic_algorithm))
- [Wikipedia: Selection](https://en.wikipedia.org/wiki/Selection_(genetic_algorithm))

