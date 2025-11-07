# Use an official Node.js runtime as a parent image
FROM oven/bun:1 as build

# Set the working directory
WORKDIR /app

# Copy parcelsData.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN bun install

# Copy the rest of the application code
COPY . .

# Build the application
RUN bun run build

# Use a lightweight web server to serve the built files
FROM nginx:alpine

# Copy the built files to the nginx html directory
COPY --from=build /app/dist /opt/build
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]