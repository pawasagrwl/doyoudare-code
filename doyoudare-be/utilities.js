import { randomBytes } from 'crypto';

// Generates a secure random string for the state
const generateRandomString = (length) => {
  return randomBytes(length).toString('hex');
};


export { generateRandomString };
