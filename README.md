# Prompt2Video - AI-Powered Educational Video Generator

Transform any concept into a 15-second animated educational video using AI and Manim.

## 🚀 Features

- **Simple Input**: Enter any concept you want explained
- **AI-Powered**: Uses LLM to generate explanations and Manim code
- **Automated Animation**: Creates educational videos using Manim
- **REST API**: Clean API for integration with frontend applications
- **Job Queue**: Async processing with Redis-backed job management
- **Docker Support**: Easy deployment with Docker Compose

## 🏗️ Architecture

```
Frontend (React) → API Gateway → Go Backend → LLM Service → Manim Renderer
                                      ↓
                              Redis Job Queue ← Workers
                                      ↓
                              PostgreSQL Database
```

## 📁 Project Structure

```
prompt2video/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   └── server.go            # HTTP server setup
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database initialization
│   ├── handlers/
│   │   └── video_handler.go     # HTTP request handlers
│   ├── middleware/
│   │   └── middleware.go        # HTTP middleware
│   ├── models/
│   │   └── models.go            # Data models
│   └── services/
│       ├── llm_service.go       # LLM integration
│       └── video_service.go     # Video rendering
├── docker-compose.yml           # Docker services
├── Dockerfile                   # API container
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
└── .env.example                 # Environment template
```

## 🛠️ Setup & Installation

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (if running locally)
- Redis (if running locally)
- OpenAI API Key or Gemini API Key

### Quick Start with Docker

1. **Clone and configure**:

   ```bash
   git clone <repository>
   cd prompt2video
   cp .env.example .env
   # Edit .env with your API keys
   ```

2. **Start services**:

   ```bash
   make docker-up
   ```

3. **Test the API**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/videos/generate \
     -H "Content-Type: application/json" \
     -d '{"prompt": "Explain Pythagorean theorem", "subject": "mathematics"}'
   ```

### Local Development

1. **Install dependencies**:

   ```bash
   make deps
   make install-tools
   ```

2. **Start database services**:

   ```bash
   docker-compose up postgres redis -d
   ```

3. **Run the API**:
   ```bash
   make run
   # Or for hot reload:
   make dev
   ```

## 📡 API Endpoints

### Generate Video

```http
POST /api/v1/videos/generate
Content-Type: application/json

{
  "prompt": "Explain photosynthesis",
  "subject": "biology"
}
```

**Response:**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "queued",
  "created_at": "2025-01-01T12:00:00Z"
}
```

### Check Job Status

```http
GET /api/v1/videos/status/{jobId}
```

**Response:**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "completed",
  "video_url": "/storage/video_123_1609459200.mp4",
  "progress": 100,
  "created_at": "2025-01-01T12:00:00Z",
  "completed_at": "2025-01-01T12:00:15Z"
}
```


### List Generated Videos

```http
GET /api/v1/videos?page=1&limit=10
```


## 🔧 Configuration

Environment variables (see `.env.example`):

| Variable         | Description                  | Default            |
| ---------------- | ---------------------------- | ------------------ |
| `PORT`           | API server port              | `8080`             |
| `DATABASE_URL`   | PostgreSQL connection string |                    |
| `REDIS_URL`      | Redis connection string      |                    |
| `OPENAI_API_KEY` | OpenAI API key               |                    |
| `GEMINI_API_KEY` | Google Gemini API key        |                    |
| `VIDEO_STORAGE`  | Video storage directory      | `./storage/videos` |
| `WORKER_COUNT`   | Number of worker processes   | `3`                |

## 🎬 Video Generation Process

1. **User submits prompt** → API creates job record
2. **Job queued** → Redis manages processing queue
3. **Worker picks job** → Updates status to "processing"
4. **LLM generates content** → Creates explanation + Manim code
5. **Manim renders video** → Produces MP4 file
6. **Job completed** → Video available for download

## 🐳 Docker Setup

The application includes a complete Docker setup:

- **PostgreSQL**: Database for job and video metadata
- **Redis**: Job queue and caching
- **API Server**: Go backend service
- **Manim Container**: For video rendering (optional)

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Benchmark tests
go test -bench=. ./...
```

## 📈 Monitoring & Logging

- Health check endpoint: `GET /health`
- Structured logging with request IDs
- Metrics for job processing times
- Rate limiting per IP address

## 🚀 Deployment

### Production Deployment

1. **Build for production**:

   ```bash
   make build-linux
   ```

2. **Deploy with Docker**:

   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. **Scale workers**:
   ```bash
   docker-compose up --scale api=3
   ```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:

- Create an issue on GitHub
- Check the documentation in `/docs`
- Review the API examples in `/examples`

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Python 3.8+](https://img.shields.io/badge/python-3.8+-blue.svg)](https://www.python.org/downloads/)
[![Manim](https://img.shields.io/badge/Powered%20by-Manim-orange.svg)](https://github.com/3b1b/manim)

## ✨ What is Scenery?

Scenery is an AI-powered video generation platform that converts natural language descriptions into professional-quality educational animations. Unlike traditional AI video generators that create static content, Scenery produces **programmatic videos** using the powerful [Manim](https://github.com/3b1b/manim) library, making every element editable and customizable.

### 🎯 Key Features

- **🗣️ Conversational Interface** - Chat with AI like you would with a colleague
- **🎨 Programmatic Generation** - Every animation is code-based and fully editable
- **📚 Educational Focus** - Built specifically for teaching and training content
- **🔧 Scene-by-Scene Editing** - Modify individual scenes without starting over
- **⚡ Professional Quality** - Powered by advanced mathematical visualization tools
- **🎪 No Technical Skills Required** - Natural language to professional animation

## 🚀 Quick Start

### Prerequisites

- Python 3.8 or higher
- Node.js 16+ (for frontend)
- FFmpeg (for video processing)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/scenery.git
   cd scenery
   ```

2. **Install Python dependencies**

   ```bash
   pip install -r requirements.txt
   ```

3. **Install Manim**

   ```bash
   pip install manim
   ```

4. **Install frontend dependencies**

   ```bash
   cd frontend
   npm install
   cd ..
   ```

5. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Add your OpenAI API key and other configurations
   ```

6. **Run the application**

   ```bash
   # Start the backend
   python app.py

   # Start the frontend (in another terminal)
   cd frontend && npm start
   ```

Visit `http://localhost:3000` to start creating!

## 💬 How It Works

1. **Describe Your Concept**: Tell Scenery what you want to teach

   ```
   "Explain the Pythagorean theorem with a visual proof showing
   how the squares relate to each other"
   ```

2. **AI Generates Script**: Scenery creates Manim code for your concept
3. **Review & Edit**: Modify scenes, timing, colors, and content

4. **Generate Video**: Watch your idea come to life as a professional animation

## 🎥 Example Prompts

- _"Create a video showing how photosynthesis works with animated molecules"_
- _"Explain sorting algorithms with visual comparisons of bubble sort vs quicksort"_
- _"Show the water cycle with animated weather patterns and transformations"_
- _"Demonstrate calculus derivatives using geometric interpretations"_

## 🏗️ Project Structure

```
scenery/
├── backend/
│   ├── app.py              # Main FastAPI application
│   ├── ai/
│   │   ├── llm_client.py   # LLM integration
│   │   └── prompt_engine.py # Prompt templates
│   ├── video/
│   │   ├── manim_generator.py # Manim code generation
│   │   └── renderer.py     # Video rendering pipeline
│   └── api/
│       └── routes.py       # API endpoints
├── frontend/
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── pages/          # Main pages
│   │   └── services/       # API services
├── templates/              # Manim code templates
├── examples/               # Example projects
└── docs/                   # Documentation
```

## 🔧 Configuration

### Environment Variables

```bash
# AI Configuration
OPENAI_API_KEY=your_openai_api_key
MODEL_NAME=gpt-4

# Video Configuration
RENDER_QUALITY=high
OUTPUT_FORMAT=mp4
MAX_DURATION=300

# Database
DATABASE_URL=postgresql://user:password@localhost/scenery
```

### Manim Settings

Scenery automatically configures Manim for optimal educational video rendering:

- **Resolution**: 1080p (configurable)
- **Frame Rate**: 30fps
- **Quality**: High
- **Background**: Customizable themes

## 🎨 Customization

### Adding Custom Templates

Create new Manim templates in the `templates/` directory:

```python
# templates/physics_template.py
from manim import *

class PhysicsScene(Scene):
    def construct(self):
        # Your custom physics animation template
        pass
```

### Custom Themes

Define visual themes in `config/themes.json`:

```json
{
  "academic": {
    "background_color": "#f8f9fa",
    "primary_color": "#2c3e50",
    "accent_color": "#3498db"
  }
}
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Documentation

- [User Guide](docs/user-guide.md) - Complete usage instructions
- [API Reference](docs/api-reference.md) - Backend API documentation
- [Developer Guide](docs/developer-guide.md) - Contributing and development setup
- [Manim Integration](docs/manim-integration.md) - How Scenery works with Manim

## 🗺️ Roadmap

- [ ] **v1.0**: Core chat interface and basic Manim generation
- [ ] **v1.1**: Advanced scene editing and timeline control
- [ ] **v1.2**: Custom template marketplace
- [ ] **v1.3**: Collaborative editing features
- [ ] **v1.4**: Voice-to-animation generation
- [ ] **v2.0**: 3D animation support

## 🐛 Known Issues

- Large animations may take significant time to render
- Complex mathematical expressions need careful prompt engineering
- Video preview requires full render (working on real-time preview)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [3Blue1Brown](https://github.com/3b1b/manim) for the incredible Manim library
- [OpenAI](https://openai.com) for GPT models that power our natural language understanding
- The educational content creator community for inspiration and feedback

## 📞 Support

- 📧 Email: support@scenery.ai
- 💬 Discord: [Join our community](https://discord.gg/scenery)
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/scenery/issues)
- 📚 Docs: [Documentation Site](https://docs.scenery.ai)

---

**Made with ❤️ for educators who want to bring their ideas to life**

_Transform your teaching. Engage your audience. Make learning stick._
