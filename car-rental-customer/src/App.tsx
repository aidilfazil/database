// App.tsx
import { useState } from 'react';
import { ChakraProvider, Box, VStack, Heading } from "@chakra-ui/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import CustomerAuth from "./components/CustomerAuth";
import CarList from "./components/CarList";

const queryClient = new QueryClient();

function App() {
  const [customerId, setCustomerId] = useState<string | null>(null);

  return (
    <ChakraProvider>
      <QueryClientProvider client={queryClient}>
        <Box maxWidth="800px" margin="auto" padding={4}>
          <VStack spacing={8}>
            <Heading>Car Rental - Customer Portal</Heading>
            {!customerId ? (
              <CustomerAuth onAuth={setCustomerId} />
            ) : (
              <CarList customerId={customerId} />
            )}
          </VStack>
        </Box>
      </QueryClientProvider>
    </ChakraProvider>
  );
}

export default App;