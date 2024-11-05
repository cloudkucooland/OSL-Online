# OSL-Online

This is a Member Management System for a very unusual community.
It probably won't be useful for you.
If it looks like something that might be useful for you, feel free to use what you like, modify it to suit your specific needs.

# Throughts

If you have a million-member organization, this will not be a good tool for you.
This is not optimized at all.
Some of the things I do are dumb. (using model.MemberID.Get() in loops instead of writing proper queries)
These design decisions make for simple code, but are slow.
This is fine given that there are fewer than 2k members in our database.
We have single-digit users-per-day.
Most of the time it is one person (me) managing the data.
Sometimes others look up another person's address.
That's it.

# Other oddities

Our way of doing things doesn't fit how other organizations relate to the world.
Most of our members are vowed members, who renew each year.
This is odd for a religious order.
Some of our members have taken life-vows.
This is odd for everything but religious orders.
We are international. Most church management software is focused on a single country.
We have "friends" who aren't members, but seem a lot like members.
We don't remove folks who forget to renew.
We mark them as "friends."

