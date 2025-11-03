# Describe flow how user get the result
- User request the number for calculate the factorial, ie GET /factorial?number=10. Server check cache and return result to user if has, if not return calculating with request success. Incase no result. Publish event number 10.
- background job workers wait and handle events. Get by batch and handle. Get the max value from the batches, update factorial_max_request_numbers if the factorial_max_request_numbers in db < value of the events batches. Maybe many workers, each event process at least 1. Do simple but make sure work well even in high through put and many workers at the same times.
- A AWS step function running from bottom up to calculate the res. Step function need to save res and current number some where to use for next passing state to next step function (next_fac = (cur_number + 1) * cur_fac). (We can change from step function to any kind of background so split logic and calling, write unit tests), Use big.Int Golang:
+ Get the current number and max (if first runing, get from DB RDS. If not first, get from state passing))
+ Get the current factorial (if first runing, get from redis/s3. If not first, get from state passing)
+ Save next factorial (number + 1) status calculating
+ Calculate next factorial next_fac = (cur_number + 1) * cur_fac
+ Save to redis if in range
+ Save to S3
+ Save to metadata status done, current_number in a transaction
+ Optional save current number and result for next cycle calling.