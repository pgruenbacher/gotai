# Rules

This is the rules

## Regions

Regions are connected by edges.
Regions can have terrain: plain, mountain.
Edges have boundary types such as wall or river.
EdgeBoundaries can enact penalties on armies crossing or attacking over boundary.

## Turns
Player actions are all performed simultaneously unless exceptional cases. 

## Armies
An army can scale a wall or fjord a river at a health penalty.
An army must be deployed before it can march.
An army cannot march into a region if its occupied by a hostile enemy.
An army can attack into a region occupied by a hostile enemy.

A defeated army can only recover in a friendly territory. 
An army loses deployed status when it is in recover mode, or when it is disbanded.

An army attacking an army in a defending state will have attacking and defending bonuses respectively for the armies.
Two armies attacking one another will have attacking bonuses. Neither army changes region position after battle.
Two armies both moving into an unocuppied region will have an attacking battle.
Battles will consist of an entire turn in themselves in which cards/actions can be played. 

An allied army can support an defending army if it borders the defending army.
An allied army can support an attacking army if it borders the enemy.


# Code Organization
Regions do not perform anything, only acted upon.
There are lots of independent components: armies, navies, cities, families.
These independent components have managers.

Orders/Actions are sent to managers which will perform the order/action. After which the manager sends back the resulting event.
Events are then packaged and sent back to all of the managers again. The managers have listener sections for specific events that may affect their components. If these events meet a specific case, then the manager will perform an action on the component and release a new event.
The new events that resulted from previous events are then sent back through all the managers again. This continues until no new events are emitted from the managers. 

## Prompter

Space - finish turn
Enter - select option


## Notes
https://www.artstation.com/artist/thrax