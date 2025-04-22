# Blog Application Backend

## Overview
This project is the backend service for a comprehensive blogging platform. It manages blog posts, user accounts, comments, and social interactions, along with asynchronous email processing for notifications like account confirmations and password resets.

## Setup Instructions
- TODO: Provide setup instructions for the development environment and running the server.

## Main Features
- User registration and authentication with JWT tokens
- CRUD operations for blog posts and comments
- Support for rich content: text, markdown, image links, tags, and slugs
- Asynchronous email sending using a message broker (Nats.io initially)
- Role-based access control with enforced permissions for Users, Moderators, and Admins
- Full-text search and efficient pagination using PostgreSQL (pg_trgm)

## Technology Stack
- **Backend:** Golang, Gin/Echo, PostgreSQL (v17), Nats.io (for asynchronous processing)
- **Tools:** Docker, Git, CI/CD pipelines
