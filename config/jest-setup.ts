import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

/**
 * Configure for React 16
 */
Enzyme.configure({ adapter: new Adapter() });
