import * as moduleAlias from 'module-alias';
const sourcePath = process.env.NODE_ENV === 'development' ? 'src' : 'build';
moduleAlias.addAliases({
  '@server': sourcePath,
  '@configs': `${sourcePath}/configs`,
  '@domains': `${sourcePath}/domains`,
  '@controllers': `${sourcePath}/controllers`,
  '@applications': `${sourcePath}/applications`,
  '@utils': `${sourcePath}/utils`,
  '@infrastructures': `${sourcePath}/infrastructures`,
  '@repositories': `${sourcePath}/repositories`,
});

import 'reflect-metadata';
import { createContainer } from '@configs/container';
import { startServer } from '@configs/express';

(async () => {
  const container = createContainer();

  startServer(container);
})();
