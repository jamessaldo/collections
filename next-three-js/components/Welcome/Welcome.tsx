import { useState } from 'react';
import { Title, Text, Anchor } from '@mantine/core';
import useStyles from './Welcome.styles';
import { AppConfig } from '../../config/settings';

export function Welcome() {
  const { classes } = useStyles();
  const [apiUrl] = useState(AppConfig.API_URL);
  return (
    <>
      <Title className={classes.title} align="center" mt={100}>
        Welcome to{' '}
        <Text inherit variant="gradient" component="span">
          Mantine
        </Text>
      </Title>
      <Text color="dimmed" align="center" size="lg" sx={{ maxWidth: 580 }} mx="auto" mt="xl">
        This starter Next.js project for web-app service includes a minimal setup for
        server side rendering, if you want to learn more on Mantine + Next.js integration follow{' '}
        <Anchor href="https://mantine.dev/guides/next/" size="lg">
          this guide
        </Anchor>
        . To get started edit index.tsx file. API_URL is set to {apiUrl}
      </Text>
    </>
  );
}
