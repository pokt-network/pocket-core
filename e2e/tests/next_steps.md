**Example: End-to-end flow for sending funds from one account to another**

1. Use CLI to send funds from acc1 --> acc2 (pocket Account send)
2. Poll (periodically) until block increases (pocket query height)
3. Use CLI to check funds on acc2 (pocket Account query)
4. Use CLI to check funds on acc1 (pocket Account query)

**The bad:**

- Falls short when it comes to end-to-end transaction
- Can only be done after 1 block has done
- Does not validate transactions

**The good:**

- Gherkin & Cucumber are good

**Future customizations:**

- Poll for block increases w/ configs:
  - Timeout on block increase
  - New step pre-configs to customize:
    - Data dir
    - Block times
