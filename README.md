# SeekNEat backend

## API documentation

The authenticated routes need a http header `Authentication: Bearer <jwt token>` (without the <>)

### Authenticated

#### **POST `/api/new_meal`**

**Description:**
Creates a new meal listing.

**Request Body:**

```json
{
  "latitude": 12.345678,
  "longitude": 98.765432,
  "available_spots": 5,
  "name": "Sarmale",
  "description": "Yummy yummy #traditional #delabunica",
  "price": "$10",
  "timestamp": 1750000000
}
```

**Response:**

```json
{
  "success": true,
  "error": ""
}
```

#### **POST `/api/join_meal`**

**Description:**
Joins the current user to a meal.

**Request Body:**

```json
{
  "meal_id": 1
}
```

**Response:**

```json
{
  "success": true,
  "error": ""
}
```

#### **POST `/api/fetch_meal`**

**Description:**
Search for meals by name or description. Returns only meals with available spots and future timestamps.

**Request Body:**

```json
{
  "query": "traditional"
}
```

**Response:**

```json
{
  "success": true,
  "result": [
    {
      "id": 1,
      "latitude": 12.345678,
      "longitude": 98.765432,
      "host_id": 123,
      "timestamp": 1750000000,
      "price": "$10",
      "name": "Sarmale",
      "description": "Yummy yummy #traditional #delabunica",
      "spots": 5,
      "occupied_spots": 2
    }
  ]
}
```


#### **POST `/api/delete_meal`**

**Description:**
Deletes a meal by ID. (The authenticated user must be the host of the meal)

**Request Body:**

```json
{
  "meal_id": 1
}
```

**Response:**

```json
{
  "success": true,
  "error": ""
}
```

#### **POST `/api/change_user`**

**Description:**
Updates the user's username, email, or password.

**Request Body:** *(all fields optional, but at least one must be provided)*

```json
{
  "username": "newUsername",
  "email": "new@example.com",
  "password": "newPassword123"
}
```

**Response:**

```json
{
  "success": true,
  "error": ""
}
```

#### **GET `/api/validate_token`**

**Description:**
Checks if the JWT token is valid.

**Response:**

```json
{
  "success": true,
  "error": ""
}
```


### Unauthenticated

#### **POST `/api/login`**

**Description:**
Authenticates a user and returns a JWT token.

**Request Body:**

```json
{
  "username": "user@example.com", // Can be both username and email
  "password": "password123"
}
```

**Response (Success):**

```json
{
  "success": true,
  "error": "",
  "token": "JWT_TOKEN_HERE"
}
```

**Response (Failure):**

```json
{
  "success": false,
  "error": "Invalid credentials",
  "token": ""
}
```

---

#### **POST `/api/register`**

**Description:**
Registers a new user.

**Request Body:**

```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securePassword123"
}
```

**Response (Success):**

```json
{
  "success": true,
  "error": ""
}
```

**Response (Failure):**

```json
{
  "success": false,
  "error": "Username already belongs to an account"
}
```
