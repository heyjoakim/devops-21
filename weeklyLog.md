# Weekly Log

## Week 01

- [X] Add version control
    - [Repo](https://github.com/heyjoakim/devops-21)
- [X] Try to develop a high-level understanding of ITU-MiniTwit.
- [X] Migrate ITU-MiniTwit to run on a modern computer running Linux
    - [X] Get python to run
    - [X] Install deps
    - [X] Recompile flag_tool
    - [X] Install sqlite browser
    - [x] Run 2to3 to convert py2 to py3
    - [x] `shellcheck` and fix `control.sh`


- [X] Share Work on GitHub
    - [Repo](https://github.com/heyjoakim/devops-21)
- [X] Prep for next week 
    - [X] Discussed branching strategy, explained below in the notes section

## Notes
We meet on mondays from 10.00 - X.X.X.X (Super agile here!!)

### This is our branching strategy
PR
Develop -> Feature

Branch out from develop into a feature / bug and then create a pr to merge back into develop. From develop releases are pushed to production (maybe one test enviornment?).


## Week 02

- [X] Choose language and technology for refractoring
    - [X] And why
- [X] Choose branching strategy
- [ ] Refractor
- [ ] Commitment guidelines?
- [ ] Implement API for simulator

### Choose language and technology

|Lang/Dev|Pros|Cons|
|---|---|---|
|Go/Gorilla   |Fast compared to other suggested frameworks [[1]](https://github.com/the-benchmarker/web-frameworks), fullstack   |Setting up env can be tricky   |
|C#/ASP.NET/Razor/Blazor   |Scalable, plenty of resources, fullstack  |Somewhat heavy framework [[1]](https://github.com/the-benchmarker/web-frameworks), not easy to make 1:1 mapping, due to different structure, Still early life  |
|JS/Angualr   |Strong community, fullstack   |Not suited for 1:1 app, Not statically type |
|JS/Vue   |Easy to get started with, lightweight   |Not suited for 1:1 app, Not statically type, Needs seperate backend |

We have chosen Go as we believe this is well suited for such task and is fast compared to other frameworks.

### Choose a branching strategy
We are discussing advantages and disadvanteges between a Git Flow and Topic/Feature workflow strategy.

|Strategy | Pros | Cons
|---|---|---|
|Git flow| Seperate releases, more "controlled", more suited for weekly release | "more" work |
|Feature workflow | Continous development, cleaner Git history, Simple, Faster deploys | Need more internal communication |

We have chosen to go with a modified Git Flow strategy as we believe this is more suited for our weekly releases. We have decided to omit the release branch, since we think that it would create unecessary overhead compared to the size of the project. Our _development_-branch will do tests once CI/CD i setup. 

![](https://i.imgur.com/ea6o39W.png)

The branch structure will therefore be as following :

- `develop` All new feature branches must check out from here into feature branches and merged back into develop. The contents of the development branch would usually reflect what is deployed to the test environment.
- `main` The production branch reflects the current deployment in production. The production branch is merged with the develop branch every time a new version deployed to production.
- `feature/{feature-name}` New features are developed on feature branches following the *feature / feature name branch* structure.
-  `hotfix\{hotfix-name}` New hot fixes are developed on seperate hot fix branches following the *hotfix / hotfix branch name*