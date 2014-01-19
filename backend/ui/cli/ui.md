# Show time entries for the day

$ quantastic-client
server: localhost:12345
username: lol
password: super

$ quantastic time
from,  to,    duration, category,             description
23:00, 07:00, 8h,       Sleep,
07:00, 07:45, 45m,      Life:Morning Routine,
07:45, 07:52, 7m ,      Plan,

# Edit time

$ quantastic time --edit
