# Broadcast

- Best Effort Broadcast
- Majority-Ack Uniform Reliable Broadcast

## Exercícios

1. O Algoritmo Lazy Reliable Broadcast satisfaz o acordo regular?

O algoritmo Lazy Reliable Broadcast é uma camada sobre o Best Effort Broadcast e faz uso do Perfect Failure Detector. Quando processo _p_ qualquer que for detectado como falho pelo detector perfeito é removido de um conjunto de processos corretos e todas as mensagens que foram enviadas por esse processo _p_ e recebidas por qualquer outro processo correto são retransmitidas via broadcast. Quando as mensagens retransmitidas forem recebidas pelos demais processos corretos, essas mensagens serão armazenadas se não tivessem sido recebidas anteriormente (quando o processo _p_ ainda era correto).

Se estes processos que estiverem recebendo essa retransmissão de mensagens também tenham detectado que o processo _p_ está falho, é feita a retransmissão das mensagens que receberam. Eventualmente todos os processos que forem detectados como falhos que tenham enviado mensagens e ao menos um processo correto tenha recebido essas mensagens, todos os demais processos corretos também terão recebido essas mensagens via retransmissão (retransmite somente se o detector identificar falha).

2. O Algoritmo All-Ack Reliable Broadcast satisfaz o acordo uniforme?

Assim como o algoritmo anterior, o algoritmo All-Ack Uniform Reliable Broadcast é uma camada sobre o Best Effort Broadcast e faz uso do Perfect Failure Detector. Toda mensagem enviada por qualquer processo fica aguardando por uma mensagem de ACK de todos os processos corretos (um conjunto atualizado pelo detector de falhas). Somente quando todos os processos corretos responderem ACK que é feito o Deliver da mensagem.

Toda nova mensagem que um processo correto receber ele retransmite via broadcast para garantir que os demais processos também a recebam mesmo se já a tiverem recebido pelo próprio autor da mensagem. Eventualmente todos os processos corretos terão o mesmo conjunto de mensagens, independente de serem de processos falhos ou não. Diferente do algoritmo anterior, que retransmitia mensagens somente quando um processo era detectado como falho.

3. O Algoritmo Majority-Ack Reliable Broadcast satisfaz o acordo uniforme?

Semelhante ao algoritmo anterior, o algoritmo Majority-Ack Uniform Reliable Broadcast também é uma camada sobre o Best Effort Broadcast, porém não utiliza nenhum detector de defeitos e não aguarda uma mensagem de ACK de todos os processos corretos. Na verdade este algoritmo não faz uso de um conjunto de processos corretos, pois não pode diferenciar processos falhos de corretos sem um detector de falhas. Para contornar isso, ele parte da hipótese de que pelo menos a maioria dos processos (de um conjunto de _N_ processos) estarão corretos. 

De forma análoga ao algoritmo anterior, os processos que estiverem corretos respodem uma mensagem de ACK para o emissor da mensagem. Se a maioria dos processos enviaram uma mensagem de ACK é o bastante para fazer o processo emissor fazer Deliver desta. Eventualmente, todos os processos que estiverem corretos terão o mesmo conjunto de mensagens, independente de serem mensagens de processos falhos ou não.
