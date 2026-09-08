const AUTH_LOGIN_URL =
  import.meta.env.VITE_AUTH_LOGIN_URL ||
  "https://idiom-identity.onrender.com/idiom/login";

const handleAuthAction = () => {
  window.location.href = AUTH_LOGIN_URL;
};

const HomePage = () => {
  return (
    <div className="landing-page">
      <header className="landing-header">
        <div className="landing-brand">
          <div className="brand-mark">I</div>
          Idiom
        </div>

        <nav className="landing-nav">
          <a href="#features">Features</a>
          <a href="#architecture">How it works</a>
          <a href="#auth">Authentication</a>
          <a href="#pricing">Pricing</a>
        </nav>

        <div className="landing-actions">
          <button className="btn btn-secondary" onClick={handleAuthAction}>
            Log in
          </button>
          <button className="btn btn-primary" onClick={handleAuthAction}>
            Get started
          </button>
        </div>
      </header>

      <main>
        <section className="hero-section">
          <div className="eyebrow">
            <span className="eyebrow-dot" />
            Identity infrastructure for product teams
          </div>

          <h1>
            <span>Authentication for the modern product stack.</span>
            <span>One platform.</span>
          </h1>

          <p>
            Bring secure sign-in, identity, and application access under one roof.
            Launch workspaces, onboard projects, and ship a polished auth experience
            without maintaining the whole stack yourself.
          </p>

          <div className="hero-actions">
            <button className="btn btn-primary" onClick={handleAuthAction}>
              Start building free
            </button>
            {/* <button className="btn btn-secondary" onClick={handleAuthAction}>
              View docs
            </button> */}
          </div>
        </section>

        <section className="preview-section">
          <div className="preview-window">
            <div className="window-bar">
              <span className="window-dot" />
              <span className="window-dot" />
              <span className="window-dot" />
            </div>

            <div className="dashboard-shell">
              <aside className="dashboard-sidebar">
                <div className="sidebar-label">WORKSPACE</div>
                <strong>Northstar Lab</strong>

                <div className="side-item active">Overview</div>
                <div className="side-item">Projects</div>
                <div className="side-item">Members</div>
                <div className="side-item">Settings</div>
              </aside>

              <div className="dashboard-content">
                <div className="dashboard-header">
                  <div>
                    <h3>Projects</h3>
                    <p>Each product gets a dedicated auth workspace.</p>
                  </div>
                  <button className="chip-btn">+ New project</button>
                </div>

                <div className="project-grid">
                  <div className="project-card">
                    <div className="project-icon">W</div>
                    <h4>Web App</h4>
                    <p>notelet.vercel.app</p>
                    <div className="project-status">
                      <span className="status-dot" />
                      Production
                    </div>
                  </div>

                  <div className="project-card">
                    <div className="project-icon">M</div>
                    <h4>Mobile App</h4>
                    <p>iOS & Android</p>
                    <div className="project-status">
                      <span className="status-dot" />
                      Production
                    </div>
                  </div>

                  <div className="project-card">
                    <div className="project-icon">A</div>
                    <h4>Admin Portal</h4>
                    <p>admin.northstar.dev</p>
                    <div className="project-status">
                      <span className="status-dot" />
                      Production
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="feature-section" id="features">
          <div className="section-heading">
            <div className="label">BUILT FOR DEVELOPERS</div>
            <h2>Everything you need to own authentication.</h2>
            <p>
              Your team ships product features. We handle the identity layer under the hood.
            </p>
          </div>

          <div className="feature-grid">
            <article className="feature-card">
              <div className="feature-icon">01</div>
              <h3>Multiple applications</h3>
              <p>
                Create independent projects for each product and keep user access isolated by app.
              </p>
            </article>

            <article className="feature-card">
              <div className="feature-icon">02</div>
              <h3>Hosted auth flows</h3>
              <p>
                Give users a polished sign-in experience with secure OAuth and branded screens.
              </p>
            </article>

            <article className="feature-card">
              <div className="feature-icon">03</div>
              <h3>OAuth built in</h3>
              <p>
                Add Google and other identity providers without rewriting callbacks and token logic.
              </p>
            </article>

            <article className="feature-card">
              <div className="feature-icon">04</div>
              <h3>Secure by default</h3>
              <p>
                Sessions, tokens, and password handling are designed for production-scale apps.
              </p>
            </article>

            <article className="feature-card">
              <div className="feature-icon">05</div>
              <h3>Developer-first APIs</h3>
              <p>
                Add identity to your stack with clear APIs and streamlined integration patterns.
              </p>
            </article>

            <article className="feature-card">
              <div className="feature-icon">06</div>
              <h3>Single dashboard</h3>
              <p>
                Manage projects, teams, and configuration in one workspace from a single control point.
              </p>
            </article>
          </div>
        </section>

        <section className="architecture-section" id="architecture">
          <div className="architecture-inner">
            <div className="section-heading light">
              <div className="label">SIMPLE ARCHITECTURE</div>
              <h2>One account. Multiple workspaces. Unlimited apps.</h2>
              <p>
                Organize identity the same way your team organizes your products and releases.
              </p>
            </div>

            <div className="flow-grid">
              <div className="flow-card">
                <div className="flow-number">01 / ACCOUNT</div>
                <h3>Your account</h3>
                <p>A developer signs up once and can belong to multiple organizations.</p>
              </div>

              <div className="arrow">→</div>

              <div className="flow-card">
                <div className="flow-number">02 / WORKSPACE</div>
                <h3>Workspace</h3>
                <p>A team gets its own isolated authentication environment.</p>
              </div>

              <div className="arrow">→</div>

              <div className="flow-card">
                <div className="flow-number">03 / PROJECT</div>
                <h3>Project = App</h3>
                <p>Every project owns its own product auth settings and users.</p>
              </div>
            </div>
          </div>
        </section>

        <section className="auth-preview-section" id="auth">
          <div className="auth-copy">
            <div className="section-heading">
              <div className="label">THE END-USER EXPERIENCE</div>
              <h2>Authentication that feels familiar.</h2>
              <p>
                Your product gets a clean, intuitive sign-in flow without feeling like a generic admin tool.
              </p>
            </div>

            <p>
              When users land on Idiom, they see a branded app experience with Google access,
              email sign-in, and a layout designed for security and clarity.
            </p>
          </div>

          {/* <div className="auth-card">
            <div className="mini-logo">I</div>
            <h3>Sign in to Northstar</h3>
            <div className="sub-copy">Continue to your product workspace</div>

            <button className="google-btn" onClick={handleAuthAction}>
              <span className="google-icon">G</span>
              Continue with Google
            </button>

            <div className="divider">OR</div>

            <input type="email" placeholder="Email address" className="auth-input" />
            <button className="primary-action" onClick={handleAuthAction}>Continue</button>

            <div className="terms-copy">
              By continuing, you agree to the Terms and Privacy Policy.
            </div>
          </div> */}
        </section>

        <section className="cta-section" id="pricing">
          <h2>Ship auth. Not authentication infrastructure.</h2>
          <p>
            Create your workspace, set up your first app, and start authenticating users in minutes.
          </p>
          <button className="btn btn-primary" onClick={handleAuthAction}>
            Create your workspace
          </button>
        </section>
      </main>

      <footer className="landing-footer">
        <div className="footer-inner">
          <span>© 2026 Idiom</span>
          <span>Identity infrastructure for modern applications.</span>
        </div>
      </footer>
    </div>
  );
};

export default HomePage;
