import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import axios from "axios";

const CustomerForm = () => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [driversLicense, setDriversLicense] = useState("");

  const createCustomerMutation = useMutation({
    mutationFn: (newCustomer: any) =>
      axios.post("/api/customers", newCustomer,{ withCredentials: true }),
    onSuccess: () => {
      setName("");
      setEmail("");
      setPhoneNumber("");
      setDriversLicense("");
      alert("Customer information submitted successfully!");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createCustomerMutation.mutate({
      name,
      email,
      phone_number: phoneNumber,
      drivers_license: driversLicense,
    });
  };

  return (
    <Box as="form" onSubmit={handleSubmit} width="100%">
      <VStack spacing={4}>
        <FormControl isRequired>
          <FormLabel>Name</FormLabel>
          <Input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Enter your name"
          />
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Email</FormLabel>
          <Input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter your email"
            type="email"
          />
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Phone Number</FormLabel>
          <Input
            value={phoneNumber}
            onChange={(e) => setPhoneNumber(e.target.value)}
            placeholder="Enter your phone number"
          />
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Driver's License</FormLabel>
          <Input
            value={driversLicense}
            onChange={(e) => setDriversLicense(e.target.value)}
            placeholder="Enter your driver's license number"
          />
        </FormControl>
        <Button type="submit" colorScheme="blue">
          Submit
        </Button>
      </VStack>
    </Box>
  );
};

export default CustomerForm;