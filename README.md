# Microsservices
[EN] Project developed at the event **Avança DEV** taught by [FullCycle](https://fullcycle.com.br/).

[BR] Projeto desenvolvido no evento **Avança DEV** ministado pela [FullCycle](https://fullcycle.com.br/).

[EN] The objective of the course is to learn about microservices, which will be developed in GO, in addition to increasing the range of known programming languages.

[BR] O objetivo do curso é aprender sobre microsserviços, que serão desenvolvidos em GO, além de aumentar a gama de linguagens de programação conhecidas.

---

[BR] Configuração Exchange e Filas:

[EN] Queues and Exchange configurations:

- [BR] Configuração 1
- [EN] Settings 1

**Exchange** > Add a new exchange >
```
Name: orders_ex
Type: direct
Durability: Durable
Auto delete: No
Internal: No
```

**Queues** >  Add a new queue
```
Type: Classic
Name: orders
Durability: durable
Auto delete: No
Arguments:
  - x-dead-letter-exchange = dlx

-> Bind
Queues > Add binding to this queue
From echange: orders_ex
```

- [BR] Configuração 2
- [EN] Settings 2

**Exchange** > Add a new exchange >
```
Name: dlx
Type: direct
Durability: Durable
Auto delete: No
Internal: No
```

**Queues** >  Add a new queue
```
Type: Classic
Name: orders_dlq
Durability: durable
Auto delete: No
Arguments:
  - x-dead-letter-exchange = orders_ex
  - x-message-ttl = 3000

-> Bind
Queues > Add binding to this queue
From echange: dlx
```
