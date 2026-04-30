import { Link } from "react-router-dom";
import { AppRouter } from "./router";
import { Button } from "../shared/components/ui/button";
import { Button as RadixButton, Card, Flex, Text } from "@radix-ui/themes";
import { FeatureShell } from "../shared/components";

export function App() {
  return (
    <FeatureShell>
      <header className="app-header">
        <div>
          <p className="eyebrow">project-starter</p>
          <h1>React + Vite starter</h1>
          <p className="hero-copy">Composable starter generated with modular layers.</p>
        </div>
        <nav className="app-nav">
          <Link to="/">Home</Link>
          <Link to="/about">About</Link>
        </nav>
      </header>
      <section className="card"><Button>shadcn-ready button</Button></section>
      <section className="card">
        <Card size="2">
          <Flex direction="column" gap="3">
            <Text size="3" weight="bold">Radix Themes ready</Text>
            <Text color="gray">Accessible UI primitives are now available in your starter.</Text>
            <div>
              <RadixButton>Radix button</RadixButton>
            </div>
          </Flex>
        </Card>
      </section>
      <AppRouter />
    </FeatureShell>
  );
}
