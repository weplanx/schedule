<?php

namespace Lumen\Extra\Contracts;

interface TokenInterface
{
    /**
     * Generate token
     * @param string $scene
     * @param string $jti
     * @param string $ack
     * @param array $symbol
     * @return \Lcobucci\JWT\Token|false
     */
    public function create(string $scene, string $jti, string $ack, array $symbol = []);

    /**
     * Get token
     * @param string $tokenString
     * @return \Lcobucci\JWT\Token
     */
    public function get(string $tokenString);

    /**
     * Verification token
     * @param string $scene
     * @param string $tokenString
     * @return \stdClass
     * @throws \Exception
     */
    public function verify(string $scene, string $tokenString);
}
