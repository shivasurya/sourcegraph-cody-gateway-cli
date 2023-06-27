# Cody Gateway CLI (un-official)

This CLI project refers test cody-gateway API from command line.

![cody-logo](./assets/cody-logo.png)

## Available APIs

- Health Check API
- Version API
- Anthropic Completion API 
- OpenAI Embedding API
- OpenAI Completion API

The Completion API has both code and chat completion modes. This project uses the chat completion mode.

## Usage

cody-gateway-cli uses flags to parses options from command line. (I wish golang has salesforce oclif.io framework to easily generate and manage cli apps)

### HealthCheck API

```bash

./main --host http://localhost:9992 --debugtoken {{TOKEN}} --verbose --healthcheckapi

✅ Health check OK!
Execution completed successfully

```

### Version API

```bash
./main --host http://localhost:9992 --debugtoken {{TOKEN}} --verbose --versionapi

Version: 0.0.0+dev
Execution completed successfully

```

### OpenAI Embeddings API

```bash
sourcegraph-cody-gateway-cli % ./main --host http://localhost:9992 --embeddingsapi true --accesstoken --mode chat
Enter Cody Gateway Access Token for 
{{TOKEN}}
Enter keyphrases to embed ✨ (supports multi-line and type --END-- to terminate input):
-> hi
-> --END--
[-0.035099167 -0.020636523 -0.015421565 -0.03990691 -0.027375247 0.021122552 -0.022002658 -0.019467426 -0.009484131 -0.013129348 0.029608354 -0.00469609 -0...]
```

### Anthropic Chat Completion API

```bash
sourcegraph-cody-gateway-cli % ./main --host http://localhost:9992 --anthropicapi true --accesstoken --mode chat
Enter Cody Gateway Access Token for 
{{TOKEN}}
🪄 🪄 🪄 🪄 🪄 🪄 
Establishing Session with Anthropic AI 🪄 ✨ (supports multi-line and type --END-- to terminate input):
🪄 🪄 🪄 🪄 🪄 🪄 
You -> Can you list the beaches nearby Ontario Toronto region ?
200
🪄 🪄 🪄 🪄 🪄 🪄 
------------
Some of the popular beaches near Toronto in Ontario are:

- Cherry Beach - Located on the shore of Lake Ontario, about 5 km from downtown Toronto. It has swimming areas, picnic spots and hiking trails.
...
------------
🪄 🪄 🪄 🪄 🪄 🪄 
You -> Can you summarize about Goderich beach ?
🪄 🪄 🪄 🪄 🪄 🪄 
------------
Goderich beach is a popular beach located on the shores of Lake Huron in the town of Goderich, Ontario, Canada. Some key highlights of the beach include:

• It is a very scenic beach with clear blue waters and sandy shores. It is considered one of the prettiest beaches on Lake Huron.
...

You -> --END--
```

### OpenAI Chat Completion API

```bash
sourcegraph-cody-gateway-cli % ./main --host http://localhost:9992 --openaiapi true --accesstoken --mode chat
Enter Cody Gateway Access Token for 
{{TOKEN}}
🪄 🪄 🪄 🪄 🪄 🪄 
Establishing Session with OpenAI GPT-3.5-Turbo 🪄 ✨ (supports multi-line and type --END-- to terminate input):
🪄 🪄 🪄 🪄 🪄 🪄 
You -> Could you please list best beaches to visit in Ontario for this summer?
🪄 🪄 🪄 🪄 🪄 🪄 
------------
As an AI language model, I do not have personal preferences. However, here are some of the best beaches to visit in Ontario for this summer:

1. Wasaga Beach
2. Sandbanks Provincial Park
3. Sauble Beach
...
------------
🪄 🪄 🪄 🪄 🪄 🪄 
You -> What about Goderich beach ?
🪄 🪄 🪄 🪄 🪄 🪄 
------------
Yes, Goderich Beach is also a great beach to visit in Ontario. It is located ...
------------
🪄 🪄 🪄 🪄 🪄 🪄 
You -> --END--
```
