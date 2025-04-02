#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>
#include <sys/wait.h>
#include <time.h>
#include <pthread.h>

#define FIFO_PATH "/tmp/slowpoke_fifo"
#define FIFO_RECOVER_PATH "/tmp/slowpoke_fifo_recover"

// Global variables for thread communication
int child_exited = 0;
int child_exit_status = 0;
long long sleep_surplus = 0;

static long long get_current_time_ns() {
    struct timespec ts;
    if (clock_gettime(CLOCK_MONOTONIC, &ts) == -1) {
        perror("clock_gettime");
        exit(EXIT_FAILURE);
    }
    return (long long)ts.tv_sec * 1000000000LL + ts.tv_nsec;
}

void printtime(void) {
    printf("time lld\n", get_current_time_ns());
}

void precise_sleep(long long sleep_ns) {
    long long start_time = get_current_time_ns();
    long long elapsed_time = 0;

    /* long long common = sleep_ns < sleep_surplus ? sleep_ns : sleep_surplus; */
    /* sleep_surplus -= common; */
    /* sleep_ns -= common; */
    while (elapsed_time < sleep_ns) {
        long long remaining_time = sleep_ns - elapsed_time;

        // Use nanosleep for high-resolution sleep
        struct timespec req = {
            .tv_sec = remaining_time / 1000000000LL,
            .tv_nsec = remaining_time % 1000000000LL
        };

        // Sleep for the remaining time
        if (nanosleep(&req, NULL) == -1) {
            if (errno == EINTR) {
                // If interrupted, recalculate elapsed time and continue
                elapsed_time = get_current_time_ns() - start_time;
                continue;
            } else {
                perror("nanosleep");
                exit(EXIT_FAILURE);
            }
        }

        // Calculate elapsed time after sleep
        elapsed_time = get_current_time_ns() - start_time;
    }
    /* sleep_surplus += elapsed_time - sleep_ns; */
    /* printf("done sleeping, surplus %lld\n", sleep_surplus); */
}

// Function for the FIFO monitoring thread
void *monitor_fifo(void *arg) {
    // int ret = nice(-15);
    // if (ret == -1 && errno != 0) {
    //     perror("nice");
    //     pthread_exit(NULL);
    // }
    // printf("nice set to %d\n", ret);
    // fflush(stdout);

    int fifo_fd = open(FIFO_PATH, O_RDONLY | O_NONBLOCK);
    int fifo_recover_fd = open(FIFO_RECOVER_PATH, O_WRONLY);
    int is_blocking = 0;
    char buffer[256];
    pid_t child_pgid = *((pid_t *)arg);
    if (fifo_fd == -1) {
        perror("open");
        pthread_exit(NULL);
    }
    long long accumulated_nano_sleep = 0;

    char send_buf[8];
    // send_buf[0] = 0;

    while (1) {
        // Attempt to read from the FIFO
        ssize_t bytes_read = read(fifo_fd, buffer, sizeof(buffer) - 1);
        if (bytes_read == -1 && (errno == EAGAIN || errno == EWOULDBLOCK)) {
            usleep(100000);
        } else if (bytes_read == 0) {
            // EOF (other end closed)
            usleep(100000);
        } else {
            // Null-terminate the message and print it
            if (!is_blocking) {
                is_blocking = 1;
                int flags = fcntl(fifo_fd, F_GETFL, 0);
                if (flags == -1) {
                    perror("fcntl(F_GETFL)");
                    pthread_exit(NULL);
                }
                if (fcntl(fifo_fd, F_SETFL, flags & ~O_NONBLOCK) == -1) {
                    perror("fcntl(F_SETFL)");
                    pthread_exit(NULL);
                }
            }
            long long nanosleep = *(long long *)(&buffer[0]);
            accumulated_nano_sleep += nanosleep;
            long long start_time = get_current_time_ns();
            if (kill(-child_pgid, SIGSTOP) == -1) {
                printf("error in stopping");
                fflush(stdout);
            }
            // precise_sleep(nanosleep);
            precise_sleep(accumulated_nano_sleep);
            long long end_time = get_current_time_ns();
            accumulated_nano_sleep -= (end_time - start_time);
            if (write(fifo_recover_fd, send_buf, sizeof(send_buf)) == -1) {
                perror("write");
                pthread_exit(NULL);
            }
            if (kill(-child_pgid, SIGCONT) == -1) {
                printf("error in conting");
                fflush(stdout);
            }
        }

        // Check if the child has exited
        if (child_exited) {
            break;
        }
    }

    close(fifo_fd);
    pthread_exit(NULL);
}

// Function for the child monitoring thread
void *monitor_child(void *arg) {
    pid_t pid = *((pid_t *)arg);

    // Wait for the child process to finish
    int status;
    waitpid(pid, &status, 0);

    if (!WIFEXITED(status)) {
        printf("Child process terminated abnormally\n");
    }

    // Signal the FIFO monitoring thread to exit
    child_exited = 1;
    child_exit_status = WEXITSTATUS(status);

    pthread_exit(NULL);
}

int main(int argc, char *argv[]) {
    // Check if there are enough arguments
    if (argc < 2) {
        fprintf(stderr, "Usage: %s <child_program> [args...]\n", argv[0]);
        exit(EXIT_FAILURE);
    }

    // Create the named FIFO
    if (mkfifo(FIFO_PATH, 0666) == -1) {
        if (errno != EEXIST) {
            perror("mkfifo");
            exit(EXIT_FAILURE);
        }
    }
    printf("FIFO created\n");

    // Create the recovery FIFO
    if (mkfifo(FIFO_RECOVER_PATH, 0666) == -1) {
        printf("FIFO recover failed to create\n");
        fflush(stdout);
        if (errno != EEXIST) {
            perror("mkfifo");
            exit(EXIT_FAILURE);
        }
    }
    printf("Recovery FIFO created\n");
    fflush(stdout);

    // Fork a child process
    pid_t pid = fork();
    if (pid == -1) {
        perror("fork");
        exit(EXIT_FAILURE);
    }

    if (pid == 0) { // Child process
        // Set the FIFO path as an environment variable
        int fd;
        if (setenv("SLOWPOKE_FIFO_PATH", FIFO_PATH, 1) == -1) {
            perror("setenv");
            exit(EXIT_FAILURE);
        }

        if (setenv("SLOWPOKE_FIFO_RECOVER_PATH", FIFO_RECOVER_PATH, 1) == -1) {
            perror("setenv");
            exit(EXIT_FAILURE);
        }

        // Execute the child program with the remaining arguments
        if (setsid() == -1) {
            perror("setsid");
            exit(EXIT_FAILURE);
        }

        execvp(argv[1], &argv[1]);
        perror("execvp"); // If execvp fails
        exit(EXIT_FAILURE);
    } else { // Parent process
        // Create threads
        pthread_t fifo_thread, child_thread;
        pid_t pgid = getpgid(pid);
        if (pthread_create(&fifo_thread, NULL, monitor_fifo, &pgid) != 0) {
            perror("pthread_create (fifo_thread)");
            exit(EXIT_FAILURE);
        }

        if (pthread_create(&child_thread, NULL, monitor_child, &pid) != 0) {
            perror("pthread_create (child_thread)");
            exit(EXIT_FAILURE);
        }

        // Wait for the child monitoring thread to finish
        pthread_join(child_thread, NULL);

        // Wait for the FIFO monitoring thread to finish
        pthread_join(fifo_thread, NULL);

        // Clean up
        unlink(FIFO_PATH); // Remove the FIFO file
    }

    return 0;
}
