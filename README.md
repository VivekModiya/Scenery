# Scenery 🎬

*The first AI that thinks like an animator*

**Scenery transforms educational ideas into stunning animated videos through simple conversation. The first AI that thinks like an animator, understanding your concepts and generating programmatic animations you can edit and refine. From math proofs to science processes, create engaging content that makes learning stick.**

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

- *"Create a video showing how photosynthesis works with animated molecules"*
- *"Explain sorting algorithms with visual comparisons of bubble sort vs quicksort"*
- *"Show the water cycle with animated weather patterns and transformations"*
- *"Demonstrate calculus derivatives using geometric interpretations"*

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

*Transform your teaching. Engage your audience. Make learning stick.*
