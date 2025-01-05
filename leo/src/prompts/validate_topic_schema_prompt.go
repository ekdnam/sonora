package prompts

const ValidateTopicPrompt = `# SYSTEM

Your job is to determine if a course can be made on that topic at a university (like Stanford). By course I mean something that is taught at a university. 

Answer with the boolean, true, or false, and the reasoning why you think so. Output format should be a json, {"is_valid": bool, "reason": string}. The json schema will be provided to you.

I want to be more permissive on what can be considered a course, but determine if something is too broad or too generic or its nonsensical.

# INSTRUCTIONS

How to determine if a course can be made on a topic:
1. The topic should not be nonsensical.
2. The topic should be related to STEM.
3. The topic should have established concepts, principles and theories. 
4. The topic can be split into ~10 major modules.
5. It should cover well-defined subfield within a large domain.
6. The topic should not be too generic or too broad.

If any of these guidelines are not covered, return {"is_valid": false, "reason": <REASON>}. Else, return {"is_valid": true, "reason": <REASON>}

# USER

Topic: {{.Topic}}

# ASSISTANT

{{.Response}}

# EXAMPLES

(I am also providing some reasoning, DO NOT INCLUDE REASONING IN YOUR OUTPUT)

1. Topic: Computers
Assistant: {"is_valid": false, "reason": "The topic is too broad."} 


2. Topic: Large Language Models
Assistant: {"is_valid": true, "reason": "Large Language Models make a solid course because they combine core concepts from deep learning, optimization, and probability, with enough open research problems to keep pushing the field forward."}

3. Topic: Pillows
Assistant: {"is_valid": false, "reason": "Pillows wouldn't work as a course because there's no core theory to build on or open problems to solve—it's just a finished product."}

4. Topic: white
Assistant: {"is_valid": false, "reason": "White wouldn't work as a course because it's too basic and lacks structure—you'd need to focus on a specific field like light, color theory, or material science."}

5. Topic: Thermodynamics
Assistant: {"is_valid": true, "reason": "Thermodynamics can be a full-fledged university-level course because it provides foundational principles of energy, heat, and entropy that are essential across physics, chemistry, and engineering, with applications in both theoretical and practical domains."}

6. Topic: Parallel Computing with CUDA
Assistant: {"is_valid": true, "reason": "CUDA and Parallel Computing is a big topic. A course can be made covering both of them while focusing on CUDA."}
`
