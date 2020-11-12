# /app/services

Code in this directory should handle things that are not related to the actual business logic.
This layer should only be familiar with everything related to API models.

## Should handle 

- Input validations
- Authorization
    - Role check
- Request/Response shaping
    - Extract a value from Context and replace it's corresponding one in request/response
- Errors mapping 
    
## Shouldn't handle

- Converting DTOs
- Access DB
- Implement business logic