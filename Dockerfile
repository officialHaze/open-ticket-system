FROM alpine:latest

# Set the working directory
WORKDIR /ots

# Create system GROUP and USER, add USER to GROUP
# Give necessary perms
RUN addgroup --system appgroup && \
    adduser --system appuser && \
    adduser appuser appgroup && \
    chown -R appuser:appgroup /var/log && \
    chown -R appuser:appgroup /ots && \
    chown -R appuser:appgroup /home

# Copy the pre-built binary
COPY ./bin/ots .

# Copy config files
COPY ./settings/settings.jsonc ./settings/settings.jsonc
COPY slack_config.yml .

# Expose the port
EXPOSE 14069

# Switch to non-root user
USER appuser

CMD ["./ots", "ots-server"]