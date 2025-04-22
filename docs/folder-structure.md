/blog-app/
├── README.md
├── .gitignore
├── go.mod
├── /cmd/
└── main.go                  // Entry point for the application
│       
└── /internal/
├── /models/
│   ├── user.go              // User model (with auditing fields)
│   ├── post.go              // Post model (tags, slug, auditing fields)
│   ├── comment.go           // Comment model (auditing fields)
│   ├── token.go             // Token model (confirmation, password reset tokens)
│   └── event.go             // Event model (for asynchronous event logging)
├── /handlers/
│   ├── auth_handler.go      // HTTP handlers for authentication endpoints
│   ├── post_handler.go      // HTTP handlers for blog post endpoints
│   └── comment_handler.go   // HTTP handlers for comment endpoints
├── /services/
│   ├── auth_service.go      // Business logic for authentication & authorization
│   └── email_service.go     // Business logic for asynchronous email sending
├── /routes/
│   └── routes.go            // API routing and versioning (e.g., /api/v1)
└── /database/
└── database.go              // Database connection and initialization
