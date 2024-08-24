# TODOs

- [ ] Script to send locker expiry emails, a Go binary will run in as cron job
      in the Docker image. - For Zoe
- [ ] Cache static asset (middleware?)
- [x] Auth token invalidating
- [x] CSFR middleware
- [ ] Admin routes
  - CSFR middleware will be applied to all routes (other than auth of course)
  - [ ] Middleware admin token checker
  - [ ] Auth: `PUT /auth/admin/`
    - Sends user name and password via url form to authenticate
    - Redirects to `/admin/dash/` on success with cookies set
  - [ ] Dash: `GET /admin/dash/`
    - Query database for all registration, dump data onto an html table
      - Each cell has a form button to `DELETE /admin/api/registration`
  - [ ] Remove registration: `DELETE /admin/api/registration`
    - Form data `locker` should be sent, containing the locker ID (ie, `ELW 120`).
  - [ ] Export: `GET /admin/api/export`
    - Empty body
    - Exports the current `registration` table into a csv file, then self-email
