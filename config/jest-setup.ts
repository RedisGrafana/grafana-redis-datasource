import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

/**
 * Configure for React 16
 */
Enzyme.configure({ adapter: new Adapter() });

/**
 * Workaround for Grafana 8.0.3
 */
Object.defineProperty(global, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});
