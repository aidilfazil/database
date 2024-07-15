import { ChakraProvider, Box, VStack, Heading } from "@chakra-ui/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import CarForm from "./components/CarForm";
import CarList from "./components/CarList";

const queryClient = new QueryClient();

function App() {
  return (
    <ChakraProvider>
      <QueryClientProvider client={queryClient}>
        <Box maxWidth="800px" margin="auto" padding={4}>
          <VStack spacing={8}>
            <Heading>Car Rental - Admin Portal</Heading>
            <CarForm />
            <CarList />
          </VStack>
        </Box>
      </QueryClientProvider>
    </ChakraProvider>
  );
}

export default App;