// CustomerAuth.tsx
import React from 'react';
import { useMutation } from '@tanstack/react-query';
import { Box, Input, Button, VStack, useToast, Tabs, TabList, TabPanels, Tab, TabPanel } from '@chakra-ui/react';
import axios from 'axios';
import { useState } from 'react';

interface CustomerAuthProps {
  onAuth: (customerId: string) => void;
}

const CustomerAuth: React.FC<CustomerAuthProps> = ({ onAuth }) => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [driversLicense, setDriversLicense] = useState('');
  const toast = useToast();

  const signUpMutation = useMutation({
    mutationFn: (customerData: { name: string; email: string; phone_number: string; drivers_license: string }) =>
      axios.post('/api/customers/signup', customerData).then((res) => res.data),
    onSuccess: (data) => {
      onAuth(data.customerId);
      toast({
        title: 'Signed up successfully',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    },
    onError: (error: any) => {
      toast({
        title: 'Failed to sign up',
        description: error.response?.data?.error || 'An unknown error occurred',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    },
  });

  const loginMutation = useMutation({
    mutationFn: (customerData: { name: string; email: string }) =>
      axios.post('/api/customers/login', customerData).then((res) => res.data),
    onSuccess: (data) => {
      onAuth(data.customerId);
      toast({
        title: 'Logged in successfully',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    },
    onError: (error: any) => {
      toast({
        title: 'Failed to log in',
        description: error.response?.data?.error || 'An unknown error occurred',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    },
  });

  const handleSignUp = () => {
    signUpMutation.mutate({ name, email, phone_number: phoneNumber, drivers_license: driversLicense });
  };

  const handleLogin = () => {
    loginMutation.mutate({ name, email });
  };

  return (
    <Box>
      <Tabs>
        <TabList>
          <Tab>Sign Up</Tab>
          <Tab>Login</Tab>
        </TabList>

        <TabPanels>
          <TabPanel>
            <VStack spacing={4}>
              <Input
                placeholder="Name"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
              <Input
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <Input
                placeholder="Phone Number"
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value)}
              />
              <Input
                placeholder="Driver's License"
                value={driversLicense}
                onChange={(e) => setDriversLicense(e.target.value)}
              />
              <Button onClick={handleSignUp} isLoading={signUpMutation.isPending}>
                Sign Up
              </Button>
            </VStack>
          </TabPanel>
          <TabPanel>
            <VStack spacing={4}>
              <Input
                placeholder="Name"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
              <Input
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <Button onClick={handleLogin} isLoading={loginMutation.isPending}>
                Login
              </Button>
            </VStack>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Box>
  );
};

export default CustomerAuth;