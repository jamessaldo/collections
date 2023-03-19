import { useSelector } from 'react-redux';
import { selectAuthState, setAuthState } from '../redux/authSlice';
import { Welcome } from '../components/Welcome/Welcome';
import { ColorSchemeToggle } from '../components/ColorSchemeToggle/ColorSchemeToggle';
import { wrapper } from '../redux/store';

export default function HomePage() {
  const authState = useSelector(selectAuthState);

  return (
    <>
      <div>{authState ? 'Logged in' : 'Not Logged In'}</div>
      <Welcome />
      <ColorSchemeToggle />
    </>
  );
}

export const getServerSideProps = wrapper.getServerSideProps(
  (store) =>
    async () => {
      // we can set the initial state from here
      // we are setting to false but you can run your custom logic here
      store.dispatch(setAuthState(false));
      return {
        props: {
          authState: false,
        },
      };
    }
);
