# GTD Application Implementation Tasks

## First Sprint - Core Capture System ✅

1. **Enhance Task Model** ✅

   - ✅ Extend the existing Task model with GTD-specific fields
   - ✅ Add status transitions (inbox → other states)
   - ✅ Implement basic validation

2. **Build Capture Interface** ✅

   - ✅ Create simple input form for capturing tasks/ideas
   - ✅ Implement task list view to see captured items
   - ✅ Add basic task editing capabilities

3. **Create Basic Inbox** ✅
   - ✅ Build inbox view to show unprocessed items
   - ✅ Implement basic sorting and filtering
   - ✅ Add simple search functionality
   - ✅ Fix task views and UI consistency

## Remaining GTD Workflow Components

1. **Capture** (Advanced)

   - ✅ Implement quick-entry modal accessible from anywhere
   - ⬜ Add email integration for capturing external items

2. **Clarify/Process** ✅

   - ✅ Build inbox view for unprocessed items
   - ✅ Implement "process item" workflow with decision tree
   - ✅ Add ability to convert items to projects or next actions

3. **Organize** ✅

   - ✅ Create lists: Next Actions, Waiting For, Projects, Someday/Maybe
   - ✅ Implement contexts (location, tools, energy level)
   - ✅ Add tagging system for categorization
   - ✅ Build project management with parent/child relationships

4. **Reflect** ✅

   - ✅ Implement daily/weekly review checklists
   - ✅ Create dashboard showing upcoming deadlines
   - ⬜ Add stale item detection

5. **Engage**
   - ✅ Build context-based task filtering
   - ✅ Implement priority system for tasks
   - ⬜ Create calendar integration for time-specific commitments

## Projects Implementation

Projects are a key component of the GTD methodology for managing multi-step outcomes that require more than one action to complete. This section outlines the implementation plan for the Projects feature.

### Implementation Approach

The GTD Projects feature will be implemented using the following approach:

1. **Use existing task model**: We'll leverage the existing Task model with "project" status rather than creating a separate Project model
2. **Parent-child relationships**: Tasks can be associated with a project via the ProjectID field
3. **Progress tracking**: A project's progress will be calculated based on the completion status of its child tasks
4. **Separate interface**: Projects will have dedicated UI views while maintaining consistency with the rest of the app

5. **Projects View** ✅

   - ✅ Create dedicated Projects list page with project cards
   - ✅ Design project card with progress indicator and summary
   - ✅ Implement filtering and sorting options for projects
   - ✅ Add quick-add project button

6. **Project Detail View** ✅

   - ✅ Build project detail page showing project information
   - ✅ Implement next actions list associated with the project
   - ✅ Add project completion tracking and status indicators
   - ✅ Create project timeline/deadline visualization
   - ⬜ Add notes section for project planning

7. **Task-Project Relationship** ✅

   - ✅ Enhance task creation to allow assigning to a project
   - ✅ Create UI for moving tasks between projects
   - ✅ Implement child task creation within a project
   - ⬜ Add dependency tracking between project tasks (optional)

8. **Project Management** ✅
   - ✅ Implement project status transitions (active, on-hold, completed)
   - ✅ Add project archiving functionality
   - ⬜ Create project templates for recurring projects (optional)
   - ✅ Implement project ownership with user isolation

## Technical Tasks

1. **Data Models** ✅

   - ✅ Implement Task model with GTD-specific fields
   - ✅ Create Project model with sub-tasks (via task status)
   - ✅ Design Context and Tag models

2. **User Interface** ✅

   - ✅ Build inbox processing view
   - ✅ Create list views for different GTD categories
   - ⬜ Implement drag-and-drop for organizing tasks
   - ✅ Design weekly review interface

3. **Backend Features**

   - ✅ Implement authentication system
   - ✅ Create API endpoints for task management
   - ✅ Add recurring task functionality
   - ⬜ Build reminder system

4. **UI Enhancements**

   - ✅ Implement DaisyUI Bumblebee theme
   - ⬜ Add dark/light mode toggle
   - ⬜ Create custom CSS variables for theme consistency
   - ⬜ Implement responsive design improvements

5. **Deployment & DevOps**
   - ⬜ Set up CI/CD pipeline
   - ✅ Configure production environment (with Docker)
   - ✅ Implement backup system (via PostgreSQL integration)

## Authentication System Implementation

The authentication system will allow users to create accounts, log in, and manage their personal GTD system. This feature will support both traditional email/password authentication and social login options, starting with Google.

### 1. User Model and Database ✅

- ✅ Create `User` model with fields:
  - ID (UUID)
  - FirstName
  - LastName
  - Email (unique)
  - Password (hashed)
  - CreatedAt
  - UpdatedAt
  - LastLoginAt
  - Provider (enum: 'email', 'google', etc.)
  - ProviderID (for social login)
  - Picture (URL)
  - Active (boolean)

- ✅ Implement PostgreSQL schema and migrations
- ✅ Create database indexes for efficient querying
- ✅ Associate Tasks with User via UserID field

### 2. Authentication Core ✅

- ✅ Implement secure password hashing (using bcrypt)
- ✅ Create JWT token generation and validation
- ✅ Set up secure cookie handling
- ✅ Implement session management
- ✅ Create middleware for protected routes
- ⬜ Add CSRF protection

### 3. Traditional Authentication ✅

- ✅ Create registration page and form
- ⬜ Implement email verification flow
- ✅ Build login page with email/password
- ⬜ Add "Forgot Password" functionality
- ✅ Implement account settings page
- ⬜ Create password change functionality

### 4. Social Authentication

- ⬜ Set up OAuth 2.0 integration foundation
- ⬜ Implement Google Login
  - ⬜ Create redirect endpoints for OAuth flow
  - ⬜ Handle OAuth callbacks
  - ⬜ Extract and store user profile information
  - ⬜ Link social accounts with existing accounts
- ⬜ Add UI components for social login buttons
- ⬜ Implement user profile merging logic

### 5. Authorization & Security

- ⬜ Implement role-based access control
- ✅ Add proper error handling for authentication failures
- ✅ Create secure logout functionality
- ⬜ Implement account locking after failed attempts
- ⬜ Add IP-based suspicious activity detection
- ⬜ Create audit logging for security events

### 6. User Experience

- ✅ Design user profile page
- ⬜ Implement avatar/profile picture handling
- ✅ Create account settings page
- ⬜ Build UI for linked accounts management
- ⬜ Add session management UI (view active sessions)
- ✅ Implement personalized dashboard

### 7. Testing & QA

- ⬜ Write unit tests for authentication logic
- ⬜ Create integration tests for auth flows
- ⬜ Perform security audit of authentication system
- ⬜ Test social login workflows
- ✅ Validate error handling and edge cases