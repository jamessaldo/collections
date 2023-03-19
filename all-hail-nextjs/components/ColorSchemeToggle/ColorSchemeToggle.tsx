import { ActionIcon, Group, useMantineColorScheme } from '@mantine/core';
import { IconSun, IconMoonStars } from '@tabler/icons';
import { useDispatch, useSelector } from 'react-redux';
import { selectAuthState, setAuthState } from '../../redux/authSlice';

export function ColorSchemeToggle() {
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();
  const authState = useSelector(selectAuthState);
  const dispatch = useDispatch();

  return (
    <Group position="center" mt="xl">
      <ActionIcon
        onClick={
          () => {
            toggleColorScheme();
            console.log('selectAuthState: ', selectAuthState);
            dispatch(setAuthState(!authState));
          }
        }
        size="xl"
        sx={(theme) => ({
          backgroundColor:
            theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
          color: theme.colorScheme === 'dark' ? theme.colors.yellow[4] : theme.colors.blue[6],
        })}
      >
        {colorScheme === 'dark' ? (
          <IconSun size={20} stroke={1.5} />
        ) : (
          <IconMoonStars size={20} stroke={1.5} />
        )}
      </ActionIcon>
    </Group>
  );
}
