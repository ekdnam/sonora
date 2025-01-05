package prompts

const ValidateTopicPrompt = `# SYSTEM

Your job is to determine if the topic the user gives is indeed a topic that is related to STEM, and if a course can be made on that topic at a university (like Stanford). By course I mean something that is taught at a university. 

Answer with the boolean, true, or false. Output format should be a json, {"response": bool}. The json schema will be provided to you.

# INSTRUCTIONS

1. If the topic is nonsensical, return {"response": false}
1. Determine if the topic is related to STEM. If it is not, return {"response": false}
2. Determine if the topic is not too generic and thus a course can be made on the topic. If a course cannot be made on the topic, return {"response": false}
3. Return {"response": true}

# USER

Topic: {{.Topic}}

# ASSISTANT

{{.Response}}

# EXAMPLES

1. Topic: Computers
Assistant: {"response": false}

2. Topic: Large Language Models
Assistant: {"response": true}

3. Topic: Pillows
Assistant: {"response": false}

4. Topic: CPU
Assistant: {"response": true}
`
