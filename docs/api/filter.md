# API Documentation: Get All Expenses

## Endpoint

- **HTTP Method and URL Path:** GET `/expense/filter`

## Description

The "Get All Expenses" API retrieves all expense records associated with a specific trip ID. This endpoint is pivotal for applications that manage travel expenses, allowing users to view a comprehensive list of incurred expenses for any given trip. This API is especially useful for financial reporting, tracking, and budget assessments after a trip has concluded. Integrating this API into your application provides an efficient way to fetch expense data that complements user-driven expense entry or reporting features.

## Request

- **Method:** GET

- **Path Parameters:** None.

- **Query Parameters:**
  - `id` (string, required): Represents the Trip ID for which the expenses need to be retrieved. This parameter is used to filter expenses specific to a particular trip.

- **Request Headers:**
  - `Authorization` (string, required): Authentication is required to access this endpoint. Include a valid bearer token as follows: `Authorization: Bearer <token>`.

- **Request Body:** None.

## Response

### Success Response

- **Status Code:** 200 OK

- **Response Body Example:**

  ```json
  [
      {
          "expense_id": "e123",
          "trip_id": "t456",
          "amount": 100.00,
          "currency": "USD",
          "description": "Dinner at main square",
          "date": "2023-04-10"
      },
      {
          "expense_id": "e789",
          "trip_id": "t456",
          "amount": 50.00,
          "currency": "USD",
          "description": "Taxi to hotel",
          "date": "2023-04-11"
      }
  ]
  ```

  **Important Fields:**
  - `expense_id`: Unique identifier of the expense.
  - `trip_id`: The ID of the trip to which this expense relates.
  - `amount`: The monetary amount of the expense.
  - `currency`: Currency code of the expense amount.
  - `description`: Short description explaining the expense.
  - `date`: Date when the expense occurred.

### Error Responses

- **400 - Bad Request:** Invalid parameters; typically occurs if the `id` parameter is missing or improperly formatted.
- **401 - Unauthorized:** Missing or invalid authentication token.
- **404 - Not Found:** No expenses found for the specified Trip ID.
- **500 - Internal Server Error:** Unexpected error on the server side.

## Business Logic

1. **Authenticate Request:**
   - Ensure that a valid bearer token is provided in the headers.

2. **Validate Parameters:**
   - Confirm the `id` query parameter exists and is correctly formatted.

3. **Query Database:**
   - Retrieve all expenses associated with the given Trip ID.

4. **Return Response:**
   - Format the data as a JSON array and return it with a 200 HTTP status.

## Error Handling

- **Detection:** Errors are detected at multiple stages, including authentication, parameter validation, and database query execution.
- **Response:** Errors are communicated back to the client using standard HTTP error codes with descriptive messages.
- **Validation Failures:** If validation of the `id` parameter fails, a 400 error is returned. Authentication issues yield a 401 response, while database errors that result in no records found return a 404 error.

## Tags

- #expense
- #finance
- #reporting
- #travel
- #authentication